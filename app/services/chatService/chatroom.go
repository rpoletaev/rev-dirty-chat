package chatService

import (
	"time"

	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"gopkg.in/mgo.v2/bson"
)

const archiveSize = 20

type ChatMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
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
