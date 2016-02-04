package controllers

import (
	//"fmt"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/auth"
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
		return c.NotFound("нужно пользака запилить") //c.Redirect("/account")
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
	c.Session["Authenticated"] = "true"
	c.Session["Login"] = originalAccount.Login

	if originalAccount.IsAdmin {
		c.Session["IsAdmin"] = "true"
		return c.Redirect((*App).Index)
	}

	return c.NotFound("Нужно запилить страницу для пользака")
}

func (c *Session) Drop() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}

	return c.RenderText(c.Request.Method)
	//return c.Redirect("/")
}
