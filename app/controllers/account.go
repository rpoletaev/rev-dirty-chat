package controllers

import (
	"errors"
	"fmt"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/auth"
	"regexp"
	"strings"
)

type Account struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*Account).Before, revel.BEFORE)
	revel.InterceptMethod((*Account).After, revel.AFTER)
	revel.InterceptMethod((*Account).Panic, revel.PANIC)
}

func (u *Account) Index() revel.Result {
	return u.NotFound("нужно допилить, тут будет страница с фильтром для поиска ")
}

func (u *Account) New() revel.Result {
	if !u.Authenticated() {
		return u.Render()
	} else {
		return u.RenderError(errors.New("мадам, Вы уже авторизованы!"))
	}
}

func (u *Account) Create(password, confirm_password, email, login string) revel.Result {
	var (
		loginForm models.Account
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
		return u.Redirect((*Account).New)
	}

	loginForm = models.Account{
		HashedPassword: strings.TrimSpace(password),
		Email:          strings.TrimSpace(email),
		Login:          strings.TrimSpace(login),
	}

	originalAccount, _ := auth.FindAccountByEmail(u.Services(), loginForm.Email)
	if originalAccount != nil {
		fmt.Println("orig acc: ", originalAccount)
		u.Flash.Error("Пользователь с таким Email:[%s] уже существует ", loginForm.Email)
		return u.Redirect((*Account).New)
	}

	if loginForm.Email == "losaped@gmail.com" {
		loginForm.IsAdmin = true
	}

	err := auth.InsertAccount(u.Services(), &loginForm)
	if err != nil {
		return u.RenderError(err)
	}

	// var createdAccount models.Account{}

	// createdAccount = auth.FindAccountByEmail(u.Services(), email)
	return u.Redirect("/session/new")
}

// func (u *Account) CountAccountByEmail revel.Result {
// 	u.RenderJson()
// }
