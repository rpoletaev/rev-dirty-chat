package models

import (
	"time"
)

//** TYPES

type User struct {
	AccountLogin string    `bson: "accountlogin"`
	VisibleName  string    `bson: "visiblename"`
	Sex          string    `bson: "sex"`
	Position     string    `bson: "position"`
	Interest     string    `bson: "interest"`
	DateOfBirth  time.Time `bson: "dateofbirth"`
	ShowInSearch bool      `bson: "showinsearch"`
	About        string    `bson: "about"`
	Region       string    `bson: "region"`
	Status       string    `bson: "status"`
	Avatar       string    `bson: "avatar"`
	CreateDate   time.Time `bson: "createdt"`
}

func CreateUser(account string) User {
	return User{
		AccountLogin: account,
		VisibleName:  account,
		Sex:          "man",
		Position:     "top",
		Interest:     "Укажите свои интересы",
		DateOfBirth:  time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		ShowInSearch: true,
		About:        "Что Вы можете рассказать о себе?",
		Region:       "Краснодарский край",
		Status:       "",
		Avatar:       "/public/img/avatar/noavatar.png",
		CreateDate:   time.Now(),
	}
}

func (u User) Age() int {
	return 30
}

func (u User) Zodiac() string {
	return "leo"
}

type Sex struct {
	Name    string `bson: "name"`
	Caption string `bson: "caption"`
}

type Position struct {
	Name    string `bson: "name"`
	Caption string `bson: "caption"`
}

type Region struct {
	ID   string `bson: "_id"`
	Name string `bson: "name"`
}
