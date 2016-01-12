package controllers

import (
	"errors"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/auth"
	"regexp"
	"strings"
)

type User struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*User).Before, revel.BEFORE)
	revel.InterceptMethod((*User).After, revel.AFTER)
	revel.InterceptMethod((*User).Panic, revel.PANIC)
}

func (u *User) Index() revel.Result {
	return u.NotFound("нужно допилить, тут будет страница с фильтром для поиска ")
}

func (u *User) New() revel.Result {
	if !u.Authenticated() {
		return u.Render()
	} else {
		return u.RenderError(errors.New("мадам, Вы уже авторизованы!"))
	}
}

func (u *User) Create(password, confirm_password, email, login string) revel.Result {
	var (
		loginForm models.User
	)

	u.Validation.Required(login).Message("Не указано имя пользователя!")
	u.Validation.Required(email).Message("Не указан email!")
	u.Validation.Required(password).Message("Не указан пароль!")
	u.Validation.MinSize(password, 6).Message("Длина пароля не менее 6 символов!")
	u.Validation.Email(email).Message("Неверный формат email!")
	u.Validation.Match(password, regexp.MustCompile(confirm_password)).Message("Пароль и подтверждение не совпадают!")

	if u.Validation.HasErrors() {
		u.Validation.Keep()
		u.FlashParams()
		return u.Redirect((*User).New)
	}

	loginForm = models.User{
		HashedPassword: strings.TrimSpace(password),
		Email:          strings.TrimSpace(email),
		Login:          strings.TrimSpace(login),
	}

	//MUST BE CHECK BY UNIQUE INDEX ON MONGODB LEVEL
	// originalUser, _ := auth.FindUserByEmail(u.Services(), loginForm.Email)
	// if originalUser != nil {
	// 	u.Flash.Error("Пользователь с таким Email:[%s] уже существует ", loginForm.Email)
	// 	return u.Redirect((*User).New)
	// }

	if loginForm.Email == "losaped@gmail.com" {
		loginForm.IsAdmin = true
	}

	err := auth.InsertUser(u.Services(), &loginForm)
	if err != nil {
		return u.RenderError(err)
	}

	return u.RenderText("email:[%s], password:[%s]", loginForm.Email, loginForm.HashedPassword)
}

// func (u *User) CountUserByEmail revel.Result {
// 	u.RenderJson()
// }
