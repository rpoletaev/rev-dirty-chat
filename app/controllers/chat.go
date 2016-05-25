package controllers

import (
	"fmt"
	"time"

	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/services/chatService"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
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
	rooms := map[bson.ObjectId]string{}

	//Получим Global
	globalId := "56dde7c1e4b0c05f88d03ffe"
	rooms[bson.ObjectIdHex(globalId)] = "Общий чат"

	user, err := userService.FindUserByID(c.Services(), c.Session["CurrentUserID"])
	if err == nil && user != nil {
		region, _ := chatService.GetRegionRoom(c.Services(), user.Region)
		rooms[region.ID] = region.Name
	}

	for _, v := range user.Rooms {
		rooms[v.ID] = v.Name
	}

	return c.Render(rooms)
}

func (c *Chat) CreatePrivateRoom(withUserId string) revel.Result {
	fromUser := c.Session["CurrentUserID"]
	header, err := userService.GetPrivateRoomIDWithUser(c.Services(), fromUser, withUserId)
	if err != nil {
		c.RenderError(err)
	}

	return c.Redirect(fmt.Sprintf("/chat/%s", header.ID))
}

func (c *Chat) RoomSocket(id string, ws *websocket.Conn) revel.Result {
	fmt.Println(id)
	room, err := chatService.GetRoom(c.Services(), id)
	if err != nil {
		c.RenderText(err.Error())
	}

	for !room.IsRuning {
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Println(room.IsRuning)
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

func (c *Chat) ShowPrivateRoom(userTo string) revel.Result {
	//Get room with userTo if room is missing then create it
	userFrom := c.Session["CurrentUserID"]
	header, err := userService.GetPrivateRoomIDWithUser(c.Services(), userFrom, userTo)
	if err != nil {
		c.Flash.Error("Не удалось создать комнату с пользователем %s", userTo)
		return nil
	}

	return c.Redirect(fmt.Sprintf("/chat/%s", header.ID.String()))
}

func (c *Chat) Room(id string) revel.Result {
	room, err := chatService.GetRoom(c.Services(), id)
	if err != nil {
		return c.RenderError(err)
	}

	return c.Render(room)
}
