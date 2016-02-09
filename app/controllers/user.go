package controllers

import (
	"bytes"
	"fmt"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"gopkg.in/mgo.v2/bson"
	"image"
	"reflect"
	// _ "image/jpeg"
	// _ "image/png"
)

const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
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

	return u.Render(user, positions, sexes)
	// return u.RenderJson(sexes)
}

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

func (u *User) Update(account, name, val string) revel.Result {
	if account != u.Session["Login"] {
		return u.RenderJson(struct{ Error string }{"Вы не можете выполнить это действие!"})
	}

	user := models.User{}
	uuserType := reflect.TypeOf(user)
	field, found := userType.FieldByName("ShowInSearch")
	if found {
		return u.RenderJson(struct{ Error string }{field.Type.Name()})
	} else {
		return u.RenderJson(struct{ Error string }{"type unknown"})
	}
}

func (u *User) MainImageUpload(account string, avatar []byte) revel.Result {
	if account != u.Session["Login"] {
		return u.RenderError(fmt.Errorf("Вы не можете выполнить это действие"))
	}

	u.Validation.Required(avatar).Message("файла нема")
	u.Validation.MinSize(avatar, 2*KB).
		Message("Minimum a file size of 2KB expected")
	u.Validation.MaxSize(avatar, 2*MB).
		Message("File cannot be larger than 2MB")

	// Check format of the file.
	img, format, err := image.Decode(bytes.NewReader(avatar))
	u.Validation.Required(err == nil).Key("avatar").
		Message("Incorrect file format")
	u.Validation.Required(format == "jpeg" || format == "png").Key("avatar").
		Message("JPEG or PNG file format is expected")

	// Check resolution.
	u.Validation.Required(img.Bounds().Dy() >= 150 && img.Bounds().Dx() >= 150).Key("avatar").
		Message("Minimum allowed resolution is 150x150px")

	// Handle errors.
	if u.Validation.HasErrors() {
		u.Validation.Keep()
		u.FlashParams()
		return u.RenderJson(u.Flash.Error)
	}

	small, big := helper.CreateMainImage(img, account, revel.BasePath)
	imgData := struct {
		Small string
		Big   string
	}{
		small,
		big,
	}

	findexpr := bson.M{"accountlogin": account}
	changes := bson.M{"portrait": big, "avatar": small}
	userService.UpdateUser(u.Services(), findexpr, changes)
	return u.RenderJson(imgData)
}

func (u *User) AvatarUpload(account string, avatar []byte) revel.Result {
	if account != u.Session["Login"] {
		return u.RenderError(fmt.Errorf("Вы не можете выполнить это действие"))
	}

	u.Validation.Required(avatar).Message("файла нема")
	u.Validation.MinSize(avatar, 2*KB).
		Message("Minimum a file size of 2KB expected")
	u.Validation.MaxSize(avatar, 2*MB).
		Message("File cannot be larger than 2MB")

	// Check format of the file.
	img, format, err := image.Decode(bytes.NewReader(avatar))
	u.Validation.Required(err == nil).Key("avatar").
		Message("Incorrect file format")
	u.Validation.Required(format == "jpeg" || format == "png").Key("avatar").
		Message("JPEG or PNG file format is expected")

	// Check resolution.
	u.Validation.Required(img.Bounds().Dy() >= 150 && img.Bounds().Dx() >= 150).Key("avatar").
		Message("Minimum allowed resolution is 150x150px")

	// Handle errors.
	if u.Validation.HasErrors() {
		u.Validation.Keep()
		u.FlashParams()
		return u.RenderJson(u.Flash.Error)
	}

	processedAvatar := helper.CreateAvatar(img, account, revel.BasePath)
	findexpr := bson.M{"accountlogin": account}
	changes := bson.M{"avatar": processedAvatar}
	imgData := struct {
		Small string
	}{
		processedAvatar,
	}
	userService.UpdateUser(u.Services(), findexpr, changes)
	return u.RenderJson(imgData)
}
