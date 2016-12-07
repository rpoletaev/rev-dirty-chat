package chatService

import (
	"time"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/app/services/notifyService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"github.com/rpoletaev/wskeleton"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const archiveSize = 20

type room struct {
	*models.RoomHeader
	mongo services.Service
	*wskeleton.Hub
}

func CreateRoom(header *models.RoomHeader) *room {
	return &room{
		RoomHeader: header,
		Hub:        wskeleton.CreateHub(nil),
	}
}

func (r *room) Run() {
	var mgoErr error
	r.mongo.UserId = r.ID.Hex()
	r.mongo.MongoSession, mgoErr = mongo.CopyMonotonicSession(r.mongo.UserId)
	defer r.mongo.MongoSession.Close()
	if mgoErr != nil {
		tracelog.ERROR(mgoErr, r.mongo.UserId, "Room.Run")
	}

	history, err := GetRoomArchive(&r.mongo, r.ID)
	if err != nil {
		panic(err)
	}

	archive := wskeleton.CreateArchive(archiveSize)
	for _, hs := range history {
		archive.AddBack(hs.RestoreMessage(&r.mongo))
	}
	r.Hub.SetArchive(archive)

	//Increase count of unreaded messages by user
	r.AddBeforeBroadcast(func(msg *wskeleton.Message) {
		for _, u := range r.Users {
			if msg.Data.(MessageData).User.OriginalID != u.ID.Hex() {
				notifyService.Increase(u.ID.Hex(), r.ID.Hex())
			}
		}
	})
	//Store message in mongo
	r.Hub.AddBeforeBroadcast(func(msg *wskeleton.Message) {
		go func() {
			cm := msg.Data.(MessageData)
			sm := StoredMessage{
				ID:        bson.NewObjectId(),
				CreatedAt: cm.Timestamp,
				Body:      cm.Text,
				RoomID:    r.ID,
				UserID:    bson.ObjectIdHex(cm.User.OriginalID),
			}

			err = InsertChatMessage(&r.mongo, sm)
			if err != nil {
				println(err.Error())
			}
		}()
	})
	// r.Hub.AddAfterBroadcast(func(msg *wskeleton.Message) {
	// 	notifyService.
	// })

	r.Hub.Run()
}

func (r *room) StoreMessage(userID, data string) {
	mes := StoredMessage{
		ID:        bson.NewObjectId(),
		RoomID:    r.ID,
		UserID:    bson.ObjectIdHex(userID),
		CreatedAt: int(time.Now().Unix()),
		Body:      data,
	}

	mgoerr := r.mongo.DBAction("messages", func(collection *mgo.Collection) error {
		return collection.Insert(mes)
	})

	if mgoerr != nil {
		tracelog.ALERT("Cant't save message", "aoeu", "Say", mgoerr.Error())
	}
}

type StoredMessage struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	RoomID    bson.ObjectId `bson:"roomId"`
	UserID    bson.ObjectId `bson:"userid"`
	Body      string        `bson:"body"`
	CreatedAt int           `bson:"createdAt"`
}

func (msg *StoredMessage) RestoreMessage(service *services.Service) wskeleton.Message {
	user, _ := GetChatUser(service, msg.UserID.Hex())
	return wskeleton.Message{
		"message",
		MessageData{
			*user,
			msg.CreatedAt,
			msg.Body,
			msg.RoomID.Hex(),
		},
	}
}

type MessageData struct {
	User      ChatUser
	Timestamp int
	Text      string
	RoomID    string
}

func newEvent(typ string, user ChatUser, msg, roomId string) wskeleton.Message {
	data := MessageData{user, int(time.Now().Unix()), msg, roomId}
	return wskeleton.Message{typ, data}
}
