package chatService

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"github.com/rpoletaev/wskeleton"
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
		fmt.Sprintf("/user/%s", u.AccountLogin),
	}

	//_This.userMu.Lock()
	_This.users[user.OriginalID] = user
	//_This.userMu.Unlock()
	return user
}

type chatCacheManager struct {
	//userMu sync.Mutex
	users map[string]*ChatUser

	//roomMu sync.Mutex
	rooms map[string]*room

	//regRoomMu   sync.Mutex
	regionRooms map[string]*room

	wsPeers      map[string][]*websocket.Conn
	userChannels map[string][]chan wskeleton.Message
}

var _This *chatCacheManager

func Startup() {
	_This = &chatCacheManager{
		users:        make(map[string]*ChatUser),
		rooms:        make(map[string]*room),
		regionRooms:  make(map[string]*room),
		userChannels: make(map[string][]chan wskeleton.Message),
	}
}

func GetUserSubscribe(userId string) chan wskeleton.Message {
	println("create new channel")
	channel := make(chan wskeleton.Message)
	_This.userChannels[userId] = append(_This.userChannels[userId], channel)
	return channel
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

func GetRoomIfRunning(id string) (*room, error) {
	if room, ok := _This.rooms[id]; ok {
		return room, nil
	}

	return nil, fmt.Errorf("The room is'nt found")
}

func GetRoom(service *services.Service, id string) (room *room, err error) {
	var ok bool
	if room, ok = _This.rooms[id]; ok {
		return room, nil
	}

	room, err = GetRoomByID(service, id)
	if err == nil && room != nil {
		_This.rooms[id] = room
		go room.Run()
		return room, nil
	}
	return room, fmt.Errorf("Room not found")
}

func GetRuningRoom(id string) *room {
	if room, ok := _This.rooms[id]; ok {
		return room
	}

	return nil
}

func GetUserRoomHeaders(service *services.Service, userId string) []models.RoomHeader {
	user, err := userService.FindUserByID(service, userId)
	var rooms []models.RoomHeader
	if err == nil && user != nil {
		rooms = user.Rooms
		if user.Region != "574621ad282c61b7d98bf612" { //Default Empty Region
			region, _ := GetRegionRoom(service, user.Region)
			rooms = append(rooms, *region.RoomHeader)
		}
	}

	return rooms
}

// func FindRoomsByUser(service *services.Service, userRooms []string) (rooms []*models.RoomHeader, err error) {
// 	defer helper.CatchPanic(&err, service.UserId, "FindRooms")

// 	err = service.DBAction(COLLECTION,
// 		func(collection *mgo.Collection) error {
// 			return collection.Find(bson.M{"_id": bson.M{"$in": userRooms}}).All(&rooms)
// 		})

// 	if err != nil {
// 		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindRooms")
// 		return rooms, err
// 	}
// 	fmt.Println(rooms)
// 	tracelog.COMPLETED(service.UserId, "FindRooms")
// 	return rooms, err
// }

// func FindRoomByName(service *services.Service, name string) (room *Room, err error) {
// 	defer helper.CatchPanic(&err, service.UserId, "FindRoomsByName")

// 	queryMap := bson.M{"name": name}

// 	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "FindRoomsByName", "Query : %s", mongo.ToString(queryMap))

// 	err = service.DBAction(COLLECTION,
// 		func(collection *mgo.Collection) error {
// 			return collection.Find(queryMap).One(&room)
// 		})

// 	if err != nil {
// 		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindRoomsByName")
// 		return room, err
// 	}

// 	tracelog.COMPLETED(service.UserId, "FindRoomsByName")
// 	return room, err
// }

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

func GetRegionRoom(service *services.Service, regionId string) (room *room, err error) {
	var ok bool
	if room, ok = _This.regionRooms[regionId]; ok {
		return room, nil
	}

	var header *models.RoomHeader
	header, err = FindRoomByRegion(service, regionId)
	if err != nil {
		panic(err)
	}

	room = CreateRoom(header)

	if err == nil && room != nil {
		_This.regionRooms[regionId] = room
		go room.Run()
		return room, nil
	}

	return nil, fmt.Errorf("Комната не найдена")
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

	room := CreateRoom(header)
	_This.rooms[room.ID.String()] = room
	go room.Run()

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

	room := CreateRoom(header)
	_This.rooms[room.ID.String()] = room
	go room.Run()

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

func GetRoomByID(service *services.Service, id string) (r *room, err error) {
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
			if err != nil {
				panic(err)
			}

			r = CreateRoom(header)
			return err
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetRoomByID")
		return nil, err
	}

	tracelog.COMPLETED(service.UserId, "GetRoomByID")
	return r, nil
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

func GetRoomArchive(mongo *services.Service, id bson.ObjectId) (history []StoredMessage, err error) {
	defer helper.CatchPanic(&err, mongo.UserId, "GetRoomArchive")

	history = []StoredMessage{}
	err = mongo.DBAction("messages", func(col *mgo.Collection) error {
		return col.Find(bson.M{"roomId": id}).Sort("-createdAt").Limit(archiveSize).All(&history)
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetRoomArchive")
		return nil, err
	}

	return history, nil
}

func InsertChatMessage(mongo *services.Service, message StoredMessage) (err error) {
	defer helper.CatchPanic(&err, mongo.UserId, "InsertChatMessage")

	err = mongo.DBAction("messages", func(col *mgo.Collection) error {
		return col.Insert(message)
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "InsertChatMessage")
		return err
	}

	return nil
}
