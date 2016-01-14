package userService

import (
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION = "users"

func FindUser(service *services.Service, account string) (user *models.User, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindUser")

	tracelog.STARTED(service.UserId, "FindUser")

	queryMap := bson.M{"accountlogin": account}
	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "FindUser", "Query : %s", mongo.ToString(queryMap))

	// Execute the query
	user = &models.User{}
	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Find(queryMap).One(user)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindUser")
		return user, err
	}

	tracelog.COMPLETED(service.UserId, "FindUser")
	return user, err
}

func InsertUser(service *services.Service, user *models.User) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "InsertUser")

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Insert(user)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "InsertUser")
		return err
	}

	tracelog.COMPLETED(service.UserId, "InsertUser")
	return nil
}
