package controllers

import (
	"fmt"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/chatService"
	"golang.org/x/net/websocket"
	// "github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	//"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
)

type Chat struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*Chat).Before, revel.BEFORE)
	revel.InterceptMethod((*Chat).After, revel.AFTER)
	revel.InterceptMethod((*Chat).Panic, revel.PANIC)
}

func (c *Chat) Index(user models.User) revel.Result {
	rooms, err := chatService.FindRoomsByUser(c.Services(), user.ID.String())
	if err != nil {
		return c.RenderError(err)
	}

	var global *models.Room
	global, err = chatService.FindRoomByName(c.Services(), "global")
	if err != nil {
		return c.RenderError(err)
	}

	rooms = append(rooms, global)
	return c.Render(rooms)
}

func (c *Chat) Create(from string, to []string) revel.Result {
	name := "Беседа"
	users := []string{from}
	if to != nil && len(to) == 1 {
		name = to[0]
	}

	for _, user := range to {
		users = append(users, user)
	}

	room := models.Room{Name: name, Users: users}
	chatService.InsertRoom(c.Services(), &room)
	return c.Render(room)
}

func (c *Chat) Connect() revel.Result {
	//m.HandleRequest(c.Response.Out, c.Request.Request)
	return c.RenderText("from Connect")
}

func (c *Chat) RoomSocket(user string, ws *websocket.Conn) revel.Result {
	// Join the room.
	subscription := models.Subscribe()
	defer subscription.Cancel()

	models.Join(user)
	defer models.Leave(user)

	// Send down the archive.
	for _, event := range subscription.Archive {
		if websocket.JSON.Send(ws, &event) != nil {
			fmt.Println("Хуй че!!!")
			return nil
		}
	}

	// In order to select between websocket messages and subscription events, we
	// need to stuff websocket events into a channel.
	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			//fmt.Println(msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	// Now listen for new events from either the websocket or the chatroom.
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

			// Otherwise, say something.
			models.Say(user, msg)
		}
	}
	return nil
}
