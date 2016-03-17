package chatService

import (
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION = "rooms"

func FindRoomsByUser(service *services.Service, user string) (rooms []*models.Room, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindRooms")

	tracelog.STARTED(service.UserId, "FindRooms")

	queryMap := bson.M{"users": user}

	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "FindRooms", "Query : %s", mongo.ToString(queryMap))

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Find(queryMap).All(&rooms)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindRooms")
		return rooms, err
	}

	tracelog.COMPLETED(service.UserId, "FindRooms")
	return rooms, err
}

func FindRoomByName(service *services.Service, name string) (room *models.Room, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindRoomsByName")

	tracelog.STARTED(service.UserId, "FindRoomsByName")

	queryMap := bson.M{"name": name}

	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "FindRoomsByName", "Query : %s", mongo.ToString(queryMap))

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Find(queryMap).One(&room)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindRoomsByName")
		return room, err
	}

	tracelog.COMPLETED(service.UserId, "FindRoomsByName")
	return room, err
}

// func GetByID(service *services.Service, id string) (room *models.Room, err error) {
// 	defer helper.CatchPanic(&err, service.UserId, "FindByID")

// 	tracelog.STARTED(service.UserId, "FindByID")

// 	err = service.DBAction(COLLECTION,
// 		func(collection *mgo.Collection) error {
// 			bsonId := bson.ObjectIdHex("56dde7c1e4b0c05f88d03ffe")
// 			fmt.Println("%s", bsonId)
// 			//tracelog.ALERT("export objID", helper.MAIN_GO_ROUTINE, "FindByID", "bsonId %s", bsonId)
// 			return collection.FindId(bsonId).One(room)
// 		})

// 	if err != nil {
// 		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindByID")
// 		return nil, err
// 	}

// 	tracelog.COMPLETED(service.UserId, "FindByID")
// 	return room, nil
//}

func InsertRoom(service *services.Service, room *models.Room) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "InsertRoom")

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Insert(room)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "InsertRoom")
		return err
	}

	tracelog.COMPLETED(service.UserId, "InsertRoom")
	return nil
}

func UpdateRoom(service *services.Service, findCondition map[string]interface{}, changes map[string]interface{}) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "UpdateRoom")

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			findExpr := findCondition
			change := bson.M{"$set": changes}
			return collection.Update(findExpr, change)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpdateRoom")
		return err
	}

	tracelog.COMPLETED(service.UserId, "UpdateRoom")
	return nil
}
