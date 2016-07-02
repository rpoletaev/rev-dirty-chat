package models

import (
	"fmt"
	"time"

	"html/template"

	"gopkg.in/mgo.v2/bson"
)

type Post struct {
	Text      template.HTML `bson:"text"`
	Markdown  string        `bson:"markdown"`
	Author    ShortUser     `bson:"author"`
	Rating    Rating        `bson:"rating"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
}
type Article struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Published bool          `bson:"published"`
	Tags      []string      `bson:"tags"`
	Post
}

type Comment struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Thread bson.ObjectId `bson:"thread"`
	Path   string        `bson:"path"`
}

type Tag struct {
	ID       string   `bson:"_id,omitempty"`
	Sinonyms []string `bson:"sinonyms"`
}

type ShortUser struct {
	ID    string `bson:"id`
	Login string `bson:"login"`
}

func (su ShortUser) GetLink() string {
	return fmt.Sprintf("/user/%s", su.Login)
}

type Rating struct {
	NegativeRating int `bson:"nrating"`
	PositiveRating int `bson:"prating"`
	CommonRating   int `bson:"crating"`
}

func (r *Rating) ChangeRating(value int) {
	if value < 0 {
		r.NegativeRating++
		r.CommonRating--
	} else {
		r.PositiveRating++
		r.CommonRating++
	}
}
