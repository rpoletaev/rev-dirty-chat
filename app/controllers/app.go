package controllers

import (
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
)

type App struct {
	cb.BaseController
}

func (c *App) Index() revel.Result {
	if !c.Authenticated() {
		return c.Redirect("/session/new")
	}

	return c.Render()
}
