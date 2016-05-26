package userService

import (
	"fmt"
	"time"

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

func GetPrivateRoomIDWithUser(service *services.Service, fromUser, toUser string) (*models.RoomHeader, error) {
	userRooms := []struct {
		UserID    bson.ObjectId `bson:"_id"`
		Code      bson.ObjectId `bson:"code"`
		Name      string        `bson:"name"`
		IsPrivate bool          `bson:"is_private"`
		Region    string        `bson:"region,omitempty"`
		Users     []string      `bson:"users"`
	}{}
	var header *models.RoomHeader

	pipeline := []bson.M{
		bson.M{"$unwind": "$rooms"},
		bson.M{"$project": bson.M{
			"name":    "$rooms.name",
			"code":    "$rooms._id",
			"users":   "$rooms.users",
			"private": "$rooms.is_private",
		},
		},
		bson.M{"$match": bson.M{"$and": []bson.M{
			{"users": bson.M{"$all": []interface{}{fromUser, toUser}}},
			{"private": true},
		}}},
	}

	err := service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Pipe(pipeline).One(&userRooms)
		})

	if err != nil {
		return header, err
	}

	if len(userRooms) == 2 {
		if userRooms[0].Code == userRooms[1].Code {
			//TODO: Нужно удостовериться, что 2 комнаты у обоих, а не у одного
			if userRooms[0].UserID == userRooms[1].UserID {
				return header, fmt.Errorf("Комнаты найдены только у пользователя %s", userRooms[0].UserID.String())
			}

			header = &models.RoomHeader{
				ID:        userRooms[0].Code,
				Name:      userRooms[0].Name,
				IsPrivate: true,
				Users:     userRooms[0].Users,
			}

			return header, nil
		}
	}

	if len(userRooms) == 1 {
		var tmpUser string
		if userRooms[0].UserID.String() == fromUser {
			tmpUser = toUser
		} else if userRooms[0].UserID.String() == toUser {
			tmpUser = fromUser
		}

		tempUser, err := FindUserByID(service, tmpUser)
		if err != nil {
			fmt.Println("User with id %s not found", tmpUser)
			fmt.Println(err)
			return nil, err
		}

		header = &models.RoomHeader{
			ID:         userRooms[0].Code,
			Name:       " Room between two users",
			Users:      []string{toUser, fromUser},
			IsPrivate:  true,
			CreateDate: time.Now(),
		}

		tempUser.Rooms = append(tempUser.Rooms, *header)

		UpsertUser(service, tempUser)
		return header, nil
	}

	return header, nil
	//TODO: убрать лишние приватные комнаты между этими двумя пользователями,
	//  берем первую попавшуюся или последнюю созданную, обновляем в истории
	// сообщений коды левых комнат на выбранную, после чего чистим у пользователей
	// все левые приватки
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
