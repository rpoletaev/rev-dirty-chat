package controllers

import (
	//"errors"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	// "github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
)

type User struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*User).Before, revel.BEFORE)
	revel.InterceptMethod((*User).After, revel.AFTER)
	revel.InterceptMethod((*User).Panic, revel.PANIC)
}

func (u *User) Edit(account string) revel.Result {
	if !u.Authenticated() {
		return u.Redirect("/session/new")
	}

	user, err := userService.FindUser(u.Services(), account)
	if err != nil {
		return u.NotFound("Пользователь [%s] не найден", account)
	}

	positions := userService.GetPositions()
	for i := 0; i < len(positions); i++ {
		positions[i].Current = positions[i].Name == user.Position.Name
	}

	sexes := userService.GetSexes()
	for i := 0; i < len(sexes); i++ {
		sexes[i].Current = sexes[i].Name == user.Position.Name
	}

	// return u.Render(user, positions, sexes)
	return u.RenderJson(sexes)
}

// func (u *User) Update(account string) revel.Result {

// }

func (u *User) Create(account string) revel.Result {
	user := models.CreateUser(account)
	err := userService.InsertUser(u.Services(), &user)
	if err != nil {
		return u.RenderError(err)
	}

	return u.Redirect("/user/%s", account)
}

func (u *User) Show(account string) revel.Result {
	if !u.Authenticated() {
		return u.Redirect("/session/new")
	}

	user, err := userService.FindUser(u.Services(), account)
	if err != nil {
		return u.Redirect("/user/%s/create", account)
	}

	return u.Render(user)
}
