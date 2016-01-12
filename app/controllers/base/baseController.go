package controllerBase

import (
	"github.com/revel/revel"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
)

//** TYPES

type (
	// BaseController contains common fields and behavior for all controllers
	BaseController struct {
		*revel.Controller
		services.Service
	}
)

//** INTERCEPT FUNCTIONS

// Before is called prior to the controller method
func (this *BaseController) Before() revel.Result {
	tracelog.TRACE("BEFORE BASE SESSION", "UserId[%s] Path[%s]", this.Session.Id(), this.Request.URL.Path)

	this.UserId = this.Session.Id()
	tracelog.TRACE(this.UserId, "Before", "UserId[%s] Path[%s]", this.Session.Id(), this.Request.URL.Path)

	var err error
	this.MongoSession, err = mongo.CopyMonotonicSession(this.UserId)
	if err != nil {
		tracelog.ERRORf(err, this.UserId, "Before", this.Request.URL.Path)
		return this.RenderError(err)
	}

	return nil
}

// After is called once the controller method completes
func (this *BaseController) After() revel.Result {
	defer func() {
		if this.MongoSession != nil {
			mongo.CloseSession(this.UserId, this.MongoSession)
			this.MongoSession = nil
		}
	}()

	tracelog.TRACE(this.UserId, "After", this.Request.URL.Path)
	return nil
}

// Panic is called is an panic is caught by the framework
func (this *BaseController) Panic() revel.Result {
	defer func() {
		mongo.CloseSession(this.UserId, this.MongoSession)
		this.MongoSession = nil
	}()

	tracelog.TRACE(this.UserId, "Panic", this.Request.URL.Path)
	return nil
}

//** PUBLIC FUNCTIONS

// Base returns a pointer of the BaseController type
func (this *BaseController) Base() *BaseController {
	return this
}

// Services returns a pointer to the base services
func (this *BaseController) Services() *services.Service {
	return &this.Service
}

func (this *BaseController) Authenticated() bool {
	if this.Session["Authenticated"] == "" {
		return false
	}

	authenticated := this.Session["Authenticated"]
	return authenticated == "true"
}

func (this *BaseController) IsAdmin() bool {
	if this.Session["IsAdmin"] == "" {
		return false
	}

	isAdmin := this.Session["IsAdmin"]
	return isAdmin == "true"
}
