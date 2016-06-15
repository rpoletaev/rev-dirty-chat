package controllers

import (
	"fmt"
	"time"

	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/chatService"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"golang.org/x/net/websocket"
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
	if !c.Authenticated() {
		return c.Redirect("/session/new")
	}

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
			rooms[region.ID] = *region.RoomHeader
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

func (c *Chat) RoomSocket(id string, ws *websocket.Conn) revel.Result {
	room, err := chatService.GetRoom(c.Services(), id)
	if err != nil {
		c.RenderTemplate("errors/404.html")
	}

	for !room.IsRuning {
		time.Sleep(time.Millisecond * 500)
	}

	// Join the room.
	subscription := room.Subscribe()
	fmt.Println("We are Subscribed")
	defer room.Unsubscribe(subscription)

	room.Join(c.Services(), c.Session["CurrentUserID"])
	fmt.Println("We are Joined")
	defer room.Leave(c.Services(), c.Session["CurrentUserID"])

	// Send down the archive.
	for _, event := range subscription.Archive {
		if websocket.JSON.Send(ws, &event) != nil {
			return nil
		}
	}

	// // In order to select between websocket messages and subscription events, we
	// // need to stuff websocket events into a channel.
	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	// // Now listen for new events from either the websocket or the chatroom.
	for {
		select {
		case event := <-subscription.New:
			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
				return nil
			}
		case msg, ok := <-newMessages:
			// If the channel is closed, they disconnected.
			if !ok {
				return nil
			}

			fmt.Println(msg)
			room.Say(c.Services(), c.Session["CurrentUserID"], msg)
		}
	}
	return nil
}

func (c *Chat) Room(id string) revel.Result {
	if !bson.IsObjectIdHex(id) {
		return c.RenderTemplate("errors/501.html")
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
