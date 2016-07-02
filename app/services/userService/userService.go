package userService

import (
	"fmt"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION = "users"

func FindUser(service *services.Service, account string) (user *models.User, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindUser")

	queryMap := bson.M{"accountlogin": account}

	user = &models.User{}
	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Find(queryMap).One(user)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindUser")
		return nil, err
	}

	return user, err
}

func FindUserByID(service *services.Service, id string) (user *models.User, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindUserByID")

	user = &models.User{}
	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			if !bson.IsObjectIdHex(id) {
				return fmt.Errorf("Неправильный код")
			}
			return collection.FindId(bson.ObjectIdHex(id)).One(user)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindUser")
		return nil, err
	}

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

	return nil
}

func UpdateUser(service *services.Service, findCondition map[string]interface{}, changes map[string]interface{}) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "UpdateUser")

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			findExpr := findCondition
			change := bson.M{"$set": changes}
			return collection.Update(findExpr, change)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpdateUser")
		return err
	}

	return nil
}

func UpsertUser(service *services.Service, user *models.User) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "UpsertUser")

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			ci, er := collection.UpsertId(user.ID, user)
			fmt.Println(ci)
			return er
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpsertUser")
		return err
	}

	return nil
}

func UpdateRating(service *services.Service, userId string, value int) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "UpsertUser")

	user, err := FindUserByID(service, userId)

	if value < 0 {
		user.Rating--
	} else {
		user.Rating++
	}

	err = service.DBAction(COLLECTION, func(col *mgo.Collection) error {
		upd := bson.M{"rating": user.Rating}
		return col.UpdateId(user.ID, upd)
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpsertUser")
		return err
	}

	return nil
}
