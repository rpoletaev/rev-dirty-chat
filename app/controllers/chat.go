package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/chatService"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"github.com/rpoletaev/wskeleton"
	"gopkg.in/mgo.v2/bson"
)

type Chat struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*Chat).Before, revel.BEFORE)
	revel.InterceptMethod((*Chat).After, revel.AFTER)
	revel.InterceptMethod((*Chat).Panic, revel.PANIC)
}

func (c *Chat) Index() revel.Result {
	rooms := map[bson.ObjectId]models.RoomHeader{}

	//Получим Global
	globalId := "56dde7c1e4b0c05f88d03ffe"
	rooms[bson.ObjectIdHex(globalId)] = models.RoomHeader{
		Name: "Общий чат",
	}

	user, err := userService.FindUserByID(c.Services(), c.Session["CurrentUserID"])
	if err == nil && user != nil {
		if user.Region != "574621ad282c61b7d98bf612" { //Default Empty Region
			region, _ := chatService.GetRegionRoom(c.Services(), user.Region)
			println("reg room is ", region.RoomHeader.Name)
			rooms[region.ID] = *region.RoomHeader
			fmt.Println("RegionRoom Is ", region.RoomHeader.Name)
		}
	}

	for _, v := range user.Rooms {
		rooms[v.ID] = v
	}

	return c.Render(rooms)
}

func (c *Chat) GetPrivateRoom(user, message_text string) revel.Result {
	fromUser := c.Session["CurrentUserID"]
	usrs := []string{fromUser, user}
	header, err := chatService.GetRoomBetweenUsers(c.Services(), usrs)
	if err == nil {
		return c.Redirect(fmt.Sprintf("/chat/%s", header.ID.Hex()))
	}

	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "GetPrivateRoom", "Private room not found")
	header, err = chatService.CreatePrivateRoom(c.Services(), usrs)
	if err != nil {
		return c.RenderError(err)
	}
	return c.Redirect(fmt.Sprintf("/chat/%s", header.ID.Hex()))
}

func (c *Chat) Room(id string) revel.Result {
	if !bson.IsObjectIdHex(id) {
		return c.RenderTemplate("errors/500.html")
	}
	room, err := chatService.GetRoom(c.Services(), id)
	if err != nil {
		c.RenderTemplate("errors/404.html")
	}

	if len(room.Users) > 0 {
		user, _ := userService.FindUserByID(c.Services(), c.Session["CurrentUserID"])
		for _, r := range user.Rooms {
			if r.ID == room.ID {
				room.Name = r.Name
				room.Avatar = r.Avatar
				return c.Render(room)
			}
		}
	}

	return c.Render(room)
}

func (c *Chat) Subscribe(id string, ws *websocket.Conn) revel.Result {
	client := wskeleton.NewClient(ws, func(raw []byte) wskeleton.Message {
		var msg string
		json.Unmarshal(raw, &msg)
		cu, _ := chatService.GetChatUser(c.Services(), c.Session["CurrentUserID"])
		data := chatService.MessageData{
			User:      *cu,
			Timestamp: int(time.Now().Unix()),
			Text:      msg,
		}
		return wskeleton.Message{"message", data}
	})
	rHeaders := chatService.GetUserRoomHeaders(c.Services(), c.Session["CurrentUserID"])
	for _, rh := range rHeaders {
		if room, err := chatService.GetRoomIfRunning(rh.ID.Hex()); err == nil {
			room.RegisterClient(&client)
			go client.SendMe(room.Hub)
		}
	}

	return nil
}

func (c *Chat) RoomSocket(id string, ws *websocket.Conn) revel.Result {
	defer ws.Close()
	room, err := chatService.GetRoom(c.Services(), id)
	if err != nil {
		ws.WriteMessage(websocket.CloseMessage, []byte("room not found!"))
		return nil
	}

	frontMsg := struct {
		Event string `json:"event"`
		Data  string `json:text`
	}{}
	client := wskeleton.NewClient(ws, func(raw []byte) wskeleton.Message {
		json.Unmarshal(raw, &frontMsg)
		cu, _ := chatService.GetChatUser(c.Services(), c.Session["CurrentUserID"])
		data := chatService.MessageData{*cu, int(time.Now().Unix()), frontMsg.Data, room.ID.Hex()}
		return wskeleton.Message{"message", data}
	})

	room.RegisterClient(&client)
	go client.SendMe(room.Hub)
	client.ReadMe(room.Hub)
	return nil
}
