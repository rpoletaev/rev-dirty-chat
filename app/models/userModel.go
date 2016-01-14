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
