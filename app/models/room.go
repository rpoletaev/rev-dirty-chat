package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//RoomHeader struct to store info about room
type RoomHeader struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	IsPrivate bool          `bson:"isprivate"`
	Region    string        `bson:"region,omitempty"`
	Users     []struct {
		ID     bson.ObjectId `bson:"_id"`
		Name   string        `bson:"accountlogin"`
		Avatar string        `bson:"avatar"`
	} `bson:"users"`
	Avatar     string    `bson:"avatar"`
	CreateDate time.Time `bson:"create_date,omitempty"`
}

//KeyValuePair struct to represent short info
// about various entities
type KeyValuePair struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}

//Message Represent object to db
type Message struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	RoomID    bson.ObjectId `bson:"room_id"`
	UserID    bson.ObjectId `bson:"user_id"`
	Text      string        `bson:"text"`
	URL       string        `bson:"url"`
	CreatedAt time.Time     `bson:"created_at"`
}
