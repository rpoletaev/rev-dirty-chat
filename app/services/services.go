/*
	Services provides boilerplate functionality for all services.
	Any state required by all the services is maintained here.
*/
package services

import (
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"gopkg.in/mgo.v2"
)

//** TYPES

type (
	// Services contains common fields and behavior for all services
	Service struct {
		MongoSession *mgo.Session
		UserId       string
	}
)

//** PUBLIC FUNCTIONS

// DBAction executes queries and commands against MongoDB
func (this *Service) DBAction(collectionName string, mongoCall mongo.MongoCall) (err error) {
	return mongo.Execute(this.UserId, this.MongoSession, helper.MONGO_DATABASE, collectionName, mongoCall)
}
