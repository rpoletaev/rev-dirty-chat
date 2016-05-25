package chatService

import (
	"container/list"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
)

const archiveSize = 100

type ChatMessage struct {
	Event string `json:"event"`
	Text  string `json:"data"`
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
	Text      string
}

func newEvent(typ string, user ChatUser, msg string) Event {
	data := EventData{user, int(time.Now().Unix()), msg}
	return Event{typ, data}
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
		r.publish <- newEvent("message", *chatUser, message.Text)
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

func (r *Room) Run() {
	if r.IsRuning || r == nil {
		return
	}

	fmt.Println(r)
	if _, ok := _This.rooms[r.ID.String()]; !ok {
		_This.rooms[r.ID.String()] = r
	}

	r.subscribe = make(chan (chan<- Subscription), archiveSize)
	r.unsubscribe = make(chan (<-chan Event), archiveSize)
	r.publish = make(chan Event, archiveSize)

	archive := list.New()
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
				fmt.Println(event)
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
