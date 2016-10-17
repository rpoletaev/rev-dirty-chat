package controllers

import (
	"fmt"

	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/auth"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
)

type Session struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*Session).Before, revel.BEFORE)
	revel.InterceptMethod((*Session).After, revel.AFTER)
	revel.InterceptMethod((*Session).Panic, revel.PANIC)
}

func (c *Session) New() revel.Result {
	if !c.Authenticated() {
		return c.Render()
	} else {
		if login, ok := c.Session["Login"]; ok {
			return c.Redirect(fmt.Sprintf("/user/%s/edit", login))
		}

		return c.Render()
	}
}

func (c *Session) Create(password, email string) revel.Result {
	var (
		err             error
		originalAccount *models.Account
		loginForm       models.Account
	)

	loginForm = models.Account{
		HashedPassword: password,
		Email:          email,
	}

	originalAccount, err = auth.FindAccountByEmail(c.Services(), loginForm.Email)
	if err != nil {
		c.Flash.Error("Не удалось найти пользователя с email: ", loginForm.Email)
		return c.Redirect((*Session).New)
	}

	err = auth.VerifyPassword(loginForm.HashedPassword, originalAccount)
	if err != nil {
		c.Flash.Error("Неправильно указаны данные для входа")
		return c.Redirect((*Session).New)
	}

	//Set Session variables to valid user
	var user *models.User
	user, err = userService.FindUser(c.Services(), originalAccount.Login)
	if err != nil && err.Error() == "not found" {
		u := models.CreateUser(originalAccount.ID.Hex(), originalAccount.Login)
		user = &u
		userService.InsertUser(c.Services(), user)
	}

	user, err = userService.FindUser(c.Services(), originalAccount.Login)
	helper.CreateUserFS(revel.BasePath, originalAccount.Login)

	c.Session["Authenticated"] = "true"
	c.Session["Login"] = originalAccount.Login
	c.Session["CurrentUserID"] = user.ID.Hex()
	c.Session["VilibleName"] = user.VisibleName

	if originalAccount.IsAdmin {
		c.Session["IsAdmin"] = "true"
		return c.Redirect((*App).Index)
	}

	return c.Redirect("/user/me")
}

func (c *Session) Drop() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}

	return c.RenderText(c.Request.Method)
	//return c.Redirect("/")
}
