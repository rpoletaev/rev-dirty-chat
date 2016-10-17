package chatService

import (
	"encoding/json"
	"time"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
//archiveSize = 30
)

type room struct {
	*models.RoomHeader
	mongo      services.Service
	clients    map[*client]bool
	broadcast  chan Event
	register   chan *client
	unregister chan *client

	archive *archive
}

func CreateRoom(header *models.RoomHeader) *room {
	return &room{
		RoomHeader: header,
		broadcast:  make(chan Event),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		archive:    CreateArchive(archiveSize),
	}
}

func (r *room) run() {
	var mgoErr error

	r.mongo.UserId = r.ID.Hex()
	r.mongo.MongoSession, mgoErr = mongo.CopyMonotonicSession(r.mongo.UserId)
	defer r.mongo.MongoSession.Close()
	if mgoErr != nil {
		tracelog.ERROR(mgoErr, r.mongo.UserId, "Room.Run")
	}

	for {
		select {
		case c := <-r.register:
			r.clients[c] = true
			println("Client added! new client count is ", len(r.clients))
			r.sendArchive(c)
			break

		case c := <-r.unregister:
			_, ok := r.clients[c]
			if ok {
				delete(r.clients, c)
				close(c.send)
			}
			break

		case m := <-r.broadcast:
			r.archive.Add(m)
			r.broadcastMessage(m)
			break
		}
	}
}

func (r *room) sendArchive(c *client) {
	r.archive.Each(func(message interface{}) {
		c.send <- message.(Event)
	})
}

func (r *room) broadcastMessage(message Event) {
	for c := range r.clients {
		select {
		case c.send <- message:
			break

		default:
			close(c.send)
			delete(r.clients, c)
		}
	}
}

func (r *room) ProcessMessageFromUser(userID string, message []byte) {
	cu, err := GetChatUser(&r.mongo, userID)
	if err != nil {
		tracelog.ERROR(err, r.ID.Hex(), "ProcessMessageFromUser")
		return
	}

	var msg ChatMessage
	json.Unmarshal(message, &msg)
	r.broadcast <- newEvent("message", *cu, msg.Data, r.ID.Hex())
	r.StoreMessage(userID, msg.Data)
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

func (r *room) RegisterClient(c *client) {
	r.register <- c
}

func (r *room) FillArchive() {
	history := []StoredMessage{}

	err := r.mongo.DBAction("messages", func(col *mgo.Collection) error {
		return col.Find(bson.M{"roomId": r.ID}).Sort("_id").Limit(archiveSize).All(&history)
	})

	if err != nil {
		panic(err)
	}

	for _, hs := range history {
		r.archive.Add(hs.RestoreMessageEvent(&r.mongo))
	}
}
