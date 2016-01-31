package controllers

import (
	"bytes"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"image"
	_ "image/jpeg"
	_ "image/png"
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

func (u *User) AvatarUpload(account string, avatar []byte) revel.Result {
	u.Validation.Required(avatar).Message("файла нема")
	u.Validation.MinSize(avatar, 2*KB).
		Message("Minimum a file size of 2KB expected")
	u.Validation.MaxSize(avatar, 2*MB).
		Message("File cannot be larger than 2MB")

	// Check format of the file.
	conf, format, err := image.DecodeConfig(bytes.NewReader(avatar))
	u.Validation.Required(err == nil).Key("avatar").
		Message("Incorrect file format")
	u.Validation.Required(format == "jpeg" || format == "png").Key("avatar").
		Message("JPEG or PNG file format is expected")

	// Check resolution.
	u.Validation.Required(conf.Height >= 150 && conf.Width >= 150).Key("avatar").
		Message("Minimum allowed resolution is 150x150px")

	// Handle errors.
	if u.Validation.HasErrors() {
		u.Validation.Keep()
		u.FlashParams()
		return u.RenderJson(u.Flash.Error)
	}

	return u.RenderText("File processed")
	// return u.RenderJson(FileInfo{
	// 	ContentType: u.Params.Files["avatar"][0].Header.Get("Content-Type"),
	// 	Filename:    u.Params.Files["avatar"][0].Filename,
	// 	RealFormat:  format,
	// 	Resolution:  fmt.Sprintf("%dx%d", conf.Width, conf.Height),
	// 	Size:        len(avatar),
	// 	Status:      "Successfully uploaded",
	// })
}
