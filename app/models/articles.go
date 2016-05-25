package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	ID    bson.ObjectId `bson: "_id,omitempty"`
	Title string        `bson:"title"`
}
