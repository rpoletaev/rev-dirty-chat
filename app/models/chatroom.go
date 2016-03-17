package models

import (
	"container/list"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ChatUser struct {
	Name     string
	Avatar   string
	Status   string
	Sex      string
	Position string
	Url      string
}

func CreateChatUser(u User) ChatUser {
	return ChatUser{
		u.VisibleName,
		u.Avatar,
		u.Status,
		u.Sex.Caption,
		u.Position.Caption,
		fmt.Sprintf("/user/[%s]", u.AccountLogin),
	}
}

type Room struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Users []string      `bson: "users"`
}

type ChatMessage struct {
	Event string `json:"event"`
	Text  string `json:"data"`
}

type EventData struct {
	User      string
	Timestamp int
	Text      string
}

type Event struct {
	Type string    `json:"event"`
	Data EventData `json:"data"`
}

type Subscription struct {
	Archive []Event
	New     <-chan Event
}

func (s Subscription) Cancel() {
	unsubscribe <- s.New
}

func newEvent(typ string, user string, msg string) Event {
	data := EventData{user, int(time.Now().Unix()), msg}
	return Event{typ, data}
}

func Subscribe() Subscription {
	resp := make(chan Subscription)
	subscribe <- resp
	return <-resp
}

func Join(user string) {
	publish <- newEvent("join", user, fmt.Sprintf("%s join to the room", user))
}

func Say(user string, msg string) {
	var message ChatMessage
	json.Unmarshal([]byte(msg), &message)
	publish <- newEvent("message", user, message.Text)
}

func Leave(user string) {
	publish <- newEvent("leave", user, fmt.Sprintf("%s leave the room", user))
}

const archiveSize = 10

var (
	// Send a channel here to get room events back.  It will send the entire
	// archive initially, and then new messages as they come in.
	subscribe = make(chan (chan<- Subscription), 10)
	// Send a channel here to unsubscribe.
	unsubscribe = make(chan (<-chan Event), 10)
	// Send events here to publish them.
	publish = make(chan Event, 10)
)

// This function loops forever, handling the chat room pubsub
func chatroom() {
	archive := list.New()
	subscribers := list.New()

	for {
		select {
		case ch := <-subscribe:
			var events []Event
			for e := archive.Front(); e != nil; e = e.Next() {
				events = append(events, e.Value.(Event))
			}
			subscriber := make(chan Event, 10)
			subscribers.PushBack(subscriber)
			ch <- Subscription{events, subscriber}

		case event := <-publish:
			fmt.Println(event)
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				ch.Value.(chan Event) <- event
			}
			if archive.Len() >= archiveSize {
				archive.Remove(archive.Front())
			}
			archive.PushBack(event)

		case unsub := <-unsubscribe:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				if ch.Value.(chan Event) == unsub {
					subscribers.Remove(ch)
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}
