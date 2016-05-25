package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type RoomHeader struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`
	IsPrivate  bool          `bson:"isprivate"`
	Region     string        `bson:"region,omitempty"`
	Users      []string      `bson:"users"`
	CreateDate time.Time     `bson:"create_date,omitempty"`
}
