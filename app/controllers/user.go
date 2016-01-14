package controllers

import (
	//"errors"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"time"
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
	user, err := userService.FindUser(u.Services(), account)
	if err != nil {
		return u.Redirect("/user/%s/edit", account)
	}

	return u.RenderJson(user)
}

// func (u *User) Update(account string) revel.Result {

// }

func (u *User) Create(account string) revel.Result {
	user := models.User{
		AccountLogin: account,
		VisibleName:  account,
		Sex:          "man",
		Position:     "top",
		Interest:     "Укажите свои интересы",
		DateOfBirth:  time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		ShowInSearch: true,
		About:        "Что Вы можете рассказать о себе?",
		Region:       "Краснодарский край",
		Status:       "",
		Avatar:       "/public/img/noavatar.png",
	}
	err := userService.InsertUser(u.Services(), &user)
	if err != nil {
		return u.RenderError(err)
	}

	return u.Redirect("/user/%s", account)
}

func (u *User) Show(account string) revel.Result {
	user, err := userService.FindUser(u.Services(), account)
	if err != nil {
		return u.NotFound("Пользователь не найден")
	}

	return u.RenderJson(user)
}
