package chatService

import (
	"fmt"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	COLLECTION = "rooms"
)

type ChatUser struct {
	OriginalID string
	Name       string
	Avatar     string
	Status     string
	Sex        string
	Position   string
	Url        string
}

func (cu *ChatUser) GetValuePair() models.KeyValuePair {
	return models.KeyValuePair{ID: bson.ObjectIdHex(cu.OriginalID), Name: cu.Name}
}

func createChatUser(u *models.User) *ChatUser {
	user := &ChatUser{
		u.ID.Hex(),
		u.VisibleName,
		u.Avatar,
		u.Status,
		u.Sex.Caption,
		u.Position.Caption,
		fmt.Sprintf("/user/%s", u.AccountLogin),
	}

	_This.users[user.OriginalID] = user
	return user
}

type chatCacheManager struct {
	users       map[string]*ChatUser
	rooms       map[string]*Room
	regionRooms map[string]*Room
}

var _This *chatCacheManager

func Startup() {
	_This = &chatCacheManager{
		users:       make(map[string]*ChatUser),
		rooms:       make(map[string]*Room),
		regionRooms: make(map[string]*Room),
	}
}

func GetChatUser(service *services.Service, id string) (user *ChatUser, err error) {
	var ok bool
	if user, ok = _This.users[id]; ok {
		return user, nil
	}

	var fullUser *models.User
	fullUser, err = userService.FindUserByID(service, id)
	if err == nil {
		user = createChatUser(fullUser)

	}
	return user, err
}

func DeleteUserFromCache(id string) {
	delete(_This.users, id)
}

func GetRoom(service *services.Service, id string) (room *Room, err error) {
	var ok bool
	if room, ok = _This.rooms[id]; ok {
		return room, nil
	}

	room, err = GetRoomByID(service, id)
	if err == nil && room != nil {
		_This.rooms[id] = room
		go room.Run(service)
		return room, nil
	} else {
		return room, fmt.Errorf("Room not found")
	}
}

func FindRoomsByUser(service *services.Service, userRooms []string) (rooms []*Room, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindRooms")

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Find(bson.M{"_id": bson.M{"$in": userRooms}}).All(&rooms)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindRooms")
		return rooms, err
	}
	fmt.Println(rooms)
	tracelog.COMPLETED(service.UserId, "FindRooms")
	return rooms, err
}

func FindRoomByName(service *services.Service, name string) (room *Room, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindRoomsByName")

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

func FindRoomByRegion(service *services.Service, regionId string) (room *models.RoomHeader, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindRoomByRegion")

	queryMap := bson.M{"region": regionId}

	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "FindRoomByRegion", "Query : %s", mongo.ToString(queryMap))

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Find(queryMap).One(&room)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindRoomByRegion")
		return nil, err
	}

	tracelog.COMPLETED(service.UserId, "FindRoomByRegion")
	return room, nil
}

func GetRegionRoom(service *services.Service, regionId string) (room *Room, err error) {
	var ok bool

	if room, ok = _This.regionRooms[regionId]; ok {
		return room, nil
	}

	var header *models.RoomHeader
	header, err = FindRoomByRegion(service, regionId)
	if err != nil {
		panic(err)
	}

	room = &Room{
		RoomHeader: header,
	}
	if err == nil && room != nil {
		_This.regionRooms[regionId] = room
		go room.Run(service)
		return room, nil
	} else {
		return nil, fmt.Errorf("Комната не найдена")
	}
}

// Find RoomHeader, create if not exist
func GetRoomBetweenUsers(service *services.Service, users []string) (header *models.RoomHeader, err error) {
	queryMap := bson.M{"$and": []bson.M{bson.M{"isprivate": true}, bson.M{"users.id": bson.M{"$all": users}}}}

	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "GetRoomBetweenUsers", "Query : %s", mongo.ToString(queryMap))

	header = &models.RoomHeader{}
	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Find(queryMap).One(header)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetRoomBetweenUsers")
		return nil, err
	}

	room := &Room{RoomHeader: header}
	_This.rooms[room.ID.String()] = room
	go room.Run(service)

	tracelog.COMPLETED(service.UserId, "FindRoomsByName")
	return header, err
}

func CreatePrivateRoom(service *services.Service, users []string) (*models.RoomHeader, error) {

	kpUsrs := []struct {
		ID     bson.ObjectId `bson:"_id"`
		Name   string        `bson:"accountlogin"`
		Avatar string        `bson:"avatar"`
	}{}

	usrQuery := bson.M{"_id": bson.M{"$in": []interface{}{bson.ObjectIdHex(users[0]), bson.ObjectIdHex(users[1])}}}

	err := service.DBAction("users",
		func(collection *mgo.Collection) error {
			return collection.Find(usrQuery).Select(bson.M{"accountlogin": 1, "avatar": 1}).All(&kpUsrs)
		})

	if err != nil {
		kperror := fmt.Errorf("KeyValue for users %#v not found", users)
		tracelog.COMPLETED_ERROR(kperror, helper.MAIN_GO_ROUTINE, "GetRoomBetweenUsers")
		return nil, err
	}

	tracelog.INFO(helper.MAIN_GO_ROUTINE, "KeyValue", "%#v", kpUsrs)
	if len(kpUsrs) < 2 {
		err = fmt.Errorf("Do not all users %#v will be found", kpUsrs)
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetRoomBetweenUsers")
		return nil, err
	}

	header := &models.RoomHeader{
		ID:        bson.NewObjectId(),
		Users:     kpUsrs,
		IsPrivate: true,
	}

	room := &Room{RoomHeader: header}
	_This.rooms[room.ID.String()] = room
	go room.Run(service)

	err = InsertRoom(service, header)
	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "CreatePrivateRoom")
	}

	err = service.DBAction("users", func(collection *mgo.Collection) error {
		update := bson.M{
			"$push": bson.M{"rooms": bson.M{"_id": header.ID, "name": header.Users[1].Name, "avatar": header.Users[1].Avatar}},
		}

		return collection.UpdateId(header.Users[0].ID, update)
	})

	err = service.DBAction("users", func(collection *mgo.Collection) error {
		update := bson.M{
			"$push": bson.M{"rooms": bson.M{"_id": header.ID, "name": header.Users[0].Name, "avatar": header.Users[0].Avatar}},
		}

		return collection.UpdateId(header.Users[1].ID, update)
	})
	tracelog.COMPLETED(service.UserId, "CreatePrivateRoom")

	return header, nil
}

func GetRoomByID(service *services.Service, id string) (room *Room, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindByID")

	tracelog.STARTED(service.UserId, "FindByID")

	err = service.DBAction("rooms",
		func(collection *mgo.Collection) error {
			if !bson.IsObjectIdHex(id) {
				return fmt.Errorf("Неправильный код %s", id)
			}

			header := &models.RoomHeader{}
			query := collection.FindId(bson.ObjectIdHex(id))
			err = query.One(header)

			room = &Room{RoomHeader: header}
			if err == nil {
				room.IsRuning = false
			}

			return err
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetRoomByID")
		return nil, err
	}

	tracelog.COMPLETED(service.UserId, "GetRoomByID")
	return room, nil
}

func InsertRoom(service *services.Service, room *models.RoomHeader) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "InsertRoom")

	err = service.DBAction(COLLECTION,
		func(collection *mgo.Collection) error {
			return collection.Insert(room)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "InsertRoom")
		return err
	}

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

	return nil
}
