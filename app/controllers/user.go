package controllers

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/chatService"
	"github.com/rpoletaev/rev-dirty-chat/app/services/regionService"
	"github.com/rpoletaev/rev-dirty-chat/app/services/userService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"gopkg.in/mgo.v2/bson"
	"image"
	"image/png"
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

func (u *User) Edit() revel.Result {
	if !u.Authenticated() {
		return u.Redirect("/session/new")
	}

	user, err := userService.FindUser(u.Services(), u.Session["Login"])
	if err != nil {
		return u.RenderError(fmt.Errorf("не удалось найти пользователя %s", u.Session["Login"]))
	}

	positions := models.GetPositions()
	sexes := models.GetSexes()
	orientations := models.GetOrientations()
	var regions *[]models.Region
	regions, err = regionService.GetAllRegions(u.Services())
	if err != nil {
		fmt.Println(err)
	}
	return u.Render(user, positions, sexes, orientations, regions)
}

func (u *User) Show(account string) revel.Result {
	if !u.Authenticated() {
		return u.Redirect("/session/new")
	}

	if account == "me" {
		account = u.Session["Login"]
	}

	owner := account == u.Session["Login"]

	user, err := userService.FindUser(u.Services(), account)
	if err != nil {
		return u.RenderTemplate("errors/404.html")
	}

	return u.Render(user, owner)
}

func (u *User) Update(name, val string) revel.Result {
	user := models.User{}
	changes, err := helper.GetChangesMap(user, name, val)
	if err != nil {
		return u.RenderJson(struct{ Error string }{err.Error()})
	}
	findexpr := bson.M{"accountlogin": u.Session["Login"]}

	userService.UpdateUser(u.Services(), findexpr, changes)
	return nil
}

func (u *User) MainImageUpload(account string, avatar []byte) revel.Result {
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

	link, imgerr := helper.ImgurImageUpload(avatar, "avatar")
	if imgerr != nil {
		u.Flash.Error("%s", imgerr)
		return u.RenderJson(u.Flash.Error)
	}

	resized := imaging.Resize(img, 100, 100, imaging.Gaussian)
	buf := new(bytes.Buffer)
	png.Encode(buf, resized.SubImage(resized.Bounds()))
	smallBytes := buf.Bytes()

	small, smalerr := helper.ImgurImageUpload(smallBytes, "avatar")
	if smalerr != nil {
		u.Flash.Error("%s", smalerr)
		return u.RenderJson(u.Flash.Error)
	}
	imgData := struct {
		Small string
		Big   string
	}{
		small,
		link,
	}

	findexpr := bson.M{"accountlogin": u.Session["Login"]}
	changes := bson.M{"portrait": link, "avatar": small}
	userService.UpdateUser(u.Services(), findexpr, changes)
	chatService.DeleteUserFromCache(u.Session["CurrentUserID"])
	return u.RenderJson(imgData)
}

func (u *User) AvatarUpload(avatar []byte) revel.Result {
	u.Validation.Required(avatar).Message("файла нема")
	u.Validation.MinSize(avatar, 2*KB).
		Message("Minimum a file size of 2KB expected")
	u.Validation.MaxSize(avatar, 2*MB).
		Message("Файл не может быть более 2 мегабайт")

	// Check format of the file.
	img, format, err := image.Decode(bytes.NewReader(avatar))
	u.Validation.Required(err == nil).Key("avatar").
		Message("Incorrect file format")
	u.Validation.Required(format == "jpeg" || format == "png").Key("avatar").
		Message("Разрешенный формат: JPEG или PNG")

	// Check resolution.
	u.Validation.Required(img.Bounds().Dy() >= 100 && img.Bounds().Dx() >= 100).Key("avatar").
		Message("Размер аватара должен быть не меньше 100 на 100")

	// Handle errors.
	if u.Validation.HasErrors() {
		u.Validation.Keep()
		u.FlashParams()
		return u.RenderJson(u.Flash.Error)
	}

	processedAvatar, imgerr := helper.ImgurImageUpload(avatar, "avatar")
	if imgerr != nil {
		u.Flash.Error("%s", imgerr)
		return u.RenderJson(u.Flash.Error)
	}
	findexpr := bson.M{"accountlogin": u.Session["Login"]}
	changes := bson.M{"avatar": processedAvatar}
	imgData := struct {
		Small string
	}{
		processedAvatar,
	}
	userService.UpdateUser(u.Services(), findexpr, changes)
	chatService.DeleteUserFromCache(u.Session["CurrentUserID"])
	return u.RenderJson(imgData)
}
