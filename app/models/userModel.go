package models

import (
	//"github.com/rpoletaev/rev-dirty-chat/app/services/chatService"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	DATEFORMAT = "02 January, 2006"
)

type User struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	AccountID    string        `bson:"account_id"`
	AccountLogin string        `bson:"accountlogin"`
	VisibleName  string        `bson:"visiblename"`
	Sex          Sex           `bson:"sex"`
	Position     Position      `bson:"position"`
	Orientation  Orientation   `bson:"orientation"`
	Interest     string        `bson:"interest"`
	DateOfBirth  time.Time     `bson:"dateofbirth"`
	ShowInSearch bool          `bson:"showinsearch"`
	About        string        `bson:"about"`
	Region       string        `bson:"region"`
	Status       string        `bson:"status"`
	Avatar       string        `bson:"avatar"`
	Portrait     string        `bson:"portrait"`
	CreateDate   time.Time     `bson:"createdt"`
	Rooms        []RoomHeader  `bson:"rooms"`
}

func CreateUser(accountId, account string) User {
	return User{
		AccountLogin: account,
		AccountID:    accountId,
		VisibleName:  account,
		Sex:          Sex{Name: "man", Caption: "Мужчина"},
		Position:     Position{Name: "top", Caption: "Верх"},
		Orientation:  Orientation{Name: "hetero", Caption: "Гетеро"},
		Interest:     "Укажите свои интересы",
		DateOfBirth:  time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		ShowInSearch: true,
		About:        "Что Вы можете рассказать о себе?",
		Region:       "574621ad282c61b7d98bf612",
		Status:       "",
		Avatar:       "/public/img/avatar/noavatar.png",
		Portrait:     "/public/img/avatar/noavatar.png",
		CreateDate:   time.Now(),
	}
}

func (u *User) Age() int {
	return 30
}

func (u *User) PickerBirthDate() string {
	return u.DateOfBirth.Format(DATEFORMAT)
}
func (u *User) Zodiac() string {
	return "leo"
}

type Sex struct {
	Name    string `bson: "name"`
	Caption string `bson: "caption"`
	Current bool
}

type Position struct {
	Name    string `bson: "name"`
	Caption string `bson: "caption"`
	Current bool
}

type Orientation struct {
	Name    string `bson: "name"`
	Caption string `bson: "caption"`
	Current bool
}

func GetSexes() map[string]Sex {
	sexes := map[string]Sex{
		"woman": Sex{
			Name:    "woman",
			Caption: "Женщина",
			Current: false,
		},
		"man": Sex{
			Name:    "man",
			Caption: "Мужчина",
			Current: false,
		},
		"trans": Sex{
			Name:    "trans",
			Caption: "Транс",
			Current: false,
		},
	}

	return sexes
}

func GetPositions() map[string]Position {
	positions := map[string]Position{
		"top": Position{
			Name:    "top",
			Caption: "Верх",
			Current: false,
		},
		"bottom": Position{
			Name:    "bottom",
			Caption: "Низ",
			Current: false,
		},
		"switch": Position{
			Name:    "switch",
			Caption: "Свитч",
			Current: false,
		},
	}

	return positions
}

func GetOrientations() map[string]Orientation {
	orientations := map[string]Orientation{
		"hetero": Orientation{
			Name:    "hetero",
			Caption: "Гетеро",
			Current: false,
		},
		"homo": Orientation{
			Name:    "homo",
			Caption: "Гомо",
			Current: false,
		},
		"bi": Orientation{
			Name:    "bi",
			Caption: "Би",
			Current: false,
		},
	}

	return orientations
}

type Region struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}
