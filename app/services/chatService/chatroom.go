package chatService

import (
	"container/list"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const archiveSize = 20

type ChatMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Subscription struct {
	userId  string
	Archive []Event
	New     <-chan Event
}

func (r *Room) Unsubscribe(s Subscription) {
	r.unsubscribe <- s.New
}

type Event struct {
	Type string    `json:"event"`
	Data EventData `json:"data"`
}

type EventData struct {
	User      ChatUser
	Timestamp int
	Text      string
	RoomID    string
}

func newEvent(typ string, user ChatUser, msg, roomId string) Event {
	data := EventData{user, int(time.Now().Unix()), msg, roomId}
	return Event{typ, data}
}

type StoredMessage struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	RoomID    bson.ObjectId `bson:"roomId"`
	UserID    bson.ObjectId `bson:"userid"`
	Body      string        `bson:"body"`
	CreatedAt int           `bson:"createdAt"`
}

func (msg *StoredMessage) RestoreMessageEvent(service *services.Service) Event {
	user, _ := GetChatUser(service, msg.UserID.Hex())
	// cm := &ChatMessage{}
	// err := json.Unmarshal([]byte(msg.Body), cm)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	return Event{
		"message",
		EventData{
			*user,
			msg.CreatedAt,
			msg.Body,
			msg.RoomID.Hex(),
		},
	}
}

type Room struct {
	*models.RoomHeader
	dbService   services.Service
	IsRuning    bool
	subscribe   chan (chan<- Subscription)
	unsubscribe chan (<-chan Event)
	publish     chan Event
	inRoomSubs  map[string][]*Subscription
}

func (r *Room) Join(userId string, sub *Subscription) {
	chatUser, err := GetChatUser(&r.dbService, userId)
	if err == nil && chatUser != nil {
		r.publish <- newEvent("join", *chatUser, fmt.Sprintf("%s join to the room", chatUser.Name), r.ID.Hex())
		r.inRoomSubs[userId] = append(r.inRoomSubs[userId], sub)
	}
}

func (r *Room) Say(userId string, msg string) {
	chatUser, err := GetChatUser(&r.dbService, userId)
	if err == nil && chatUser != nil {
		var message ChatMessage
		json.Unmarshal([]byte(msg), &message)
		r.publish <- newEvent("message", *chatUser, message.Data, r.ID.Hex())

		mes := StoredMessage{
			ID:        bson.NewObjectId(),
			RoomID:    r.ID,
			UserID:    bson.ObjectIdHex(userId),
			CreatedAt: int(time.Now().Unix()),
			Body:      message.Data,
		}

		mgoerr := r.dbService.DBAction("messages", func(collection *mgo.Collection) error {
			return collection.Insert(mes)
		})

		if mgoerr != nil {
			tracelog.ALERT("Cant't save message", "aoeu", "Say", mgoerr.Error())
		}
	}
}

func (r *Room) Leave(userId string) {
	chatUser, err := GetChatUser(&r.dbService, userId)
	if err == nil && chatUser != nil {
		r.publish <- newEvent("leave", *chatUser, fmt.Sprintf("%s leave the room", chatUser.Name), r.ID.Hex())
	}
}

func (r *Room) Subscribe() Subscription {
	resp := make(chan Subscription)
	r.subscribe <- resp
	return <-resp
}

func (r *Room) Run() {
	if r.IsRuning || r == nil {
		return
	}

	var mongoErr error
	r.dbService.UserId = r.ID.Hex()
	r.dbService.MongoSession, mongoErr = mongo.CopyMonotonicSession(r.dbService.UserId)
	defer r.dbService.MongoSession.Close()
	if mongoErr != nil {
		tracelog.ERROR(mongoErr, r.ID.Hex(), "Room.Run")
	}

	r.inRoomSubs = make(map[string][]*Subscription)
	r.subscribe = make(chan (chan<- Subscription), archiveSize)
	r.unsubscribe = make(chan (<-chan Event), archiveSize)
	r.publish = make(chan Event, archiveSize)

	archive := list.New()
	history := []StoredMessage{}

	err := r.dbService.DBAction("messages", func(col *mgo.Collection) error {
		return col.Find(bson.M{"roomId": r.ID}).Sort("_id").Limit(archiveSize).All(&history)
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, hs := range history {
		archive.PushBack(hs.RestoreMessageEvent(&r.dbService))
	}
	subscribers := list.New()
	r.IsRuning = true

	for {
		select {
		case ch := <-r.subscribe:
			var events []Event
			for e := archive.Front(); e != nil; e = e.Next() {
				events = append(events, e.Value.(Event))
			}
			subscriber := make(chan Event, archiveSize)
			subscribers.PushBack(subscriber)
			ch <- Subscription{events, subscriber}

		case event := <-r.publish:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				ch.Value.(chan Event) <- event
			}
			if archive.Len() >= archiveSize {
				archive.Remove(archive.Front())
			}
			archive.PushBack(event)

		case unsub := <-r.unsubscribe:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				if ch.Value.(chan Event) == unsub {
					subscribers.Remove(ch)
					break
				}
			}
		}
	}
}
