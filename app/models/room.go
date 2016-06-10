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
