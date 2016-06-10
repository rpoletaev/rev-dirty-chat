package chatService

import (
	"container/list"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const archiveSize = 100

type ChatMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Subscription struct {
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
	Body      json.RawMessage
}

func newEvent(typ string, user ChatUser, msg json.RawMessage) Event {
	data := EventData{user, int(time.Now().Unix()), msg}
	return Event{typ, data}
}

type StoredMessage struct {
	ID        bson.ObjectId   `bson:"_id,omitempty"`
	RoomID    bson.ObjectId   `bson:"roomId"`
	UserID    bson.ObjectId   `bson: "userid"`
	Body      json.RawMessage `bson:"body"`
	CreatedAt int             `bson:"createdAt"`
}

func RestoreMessageEvent(msg StoredMessage, service *services.Service) Event {
	user, _ := GetChatUser(service, msg.UserID.Hex())
	return Event{
		"message",
		EventData{
			*user,
			msg.CreatedAt,
			msg.Body,
		},
	}
}

type Room struct {
	*models.RoomHeader
	IsRuning    bool
	subscribe   chan (chan<- Subscription)
	unsubscribe chan (<-chan Event)
	publish     chan Event
}

func (r *Room) Join(service *services.Service, userId string) {
	chatUser, err := GetChatUser(service, userId)
	if err == nil && chatUser != nil {
		r.publish <- newEvent("join", *chatUser, fmt.Sprintf("%s join to the room", chatUser.Name))
	}
}

func (r *Room) Say(service *services.Service, userId string, msg string) {
	chatUser, err := GetChatUser(service, userId)
	if err == nil && chatUser != nil {
		var message ChatMessage
		json.Unmarshal([]byte(msg), &message)
		r.publish <- newEvent("message", *chatUser, message.Data)

		mes := StoredMessage{
			ID:        bson.NewObjectId(),
			RoomID:    r.ID,
			UserID:    bson.ObjectIdHex(userId),
			CreatedAt: int(time.Now().Unix()),
			Body:      []byte(msg),
		}

		mgoerr := service.DBAction("messages", func(collection *mgo.Collection) error {
			return collection.Insert(mes)
		})

		if mgoerr != nil {
			tracelog.ALERT("Cant't save message", "aoeu", "Say", mgoerr.Error())
		}
	}
}

func (r *Room) Leave(service *services.Service, userId string) {
	chatUser, err := GetChatUser(service, userId)
	if err == nil && chatUser != nil {
		r.publish <- newEvent("leave", *chatUser, fmt.Sprintf("%s leave the room", chatUser.Name))
	}
}

func (r *Room) Subscribe() Subscription {
	resp := make(chan Subscription)
	r.subscribe <- resp
	return <-resp
}

func (r *Room) Run(service *services.Service) {
	if r.IsRuning || r == nil {
		return
	}

	r.subscribe = make(chan (chan<- Subscription), archiveSize)
	r.unsubscribe = make(chan (<-chan Event), archiveSize)
	r.publish = make(chan Event, archiveSize)

	archive := list.New()
	history := []StoredMessage{}

	err := service.DBAction("messages", func(col *mgo.Collection) error {
		return col.Find(bson.M{"roomId": r.ID}).Sort("-_id").Limit(archiveSize).All(&history)
	})

	if err == nil {
		for _, hm := range history {
			hm
		}
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

	fmt.Println("room is runing")
}
