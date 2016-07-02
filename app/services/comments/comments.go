package comments

import (
	"fmt"
	"strings"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func AddComment(service *services.Service, comment models.Comment) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "AddComment")

	tracelog.STARTED(service.UserId, "AddComment")

	err = service.DBAction("comments",
		func(collection *mgo.Collection) error {
			return collection.Insert(comment)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "AddComment")
		return err
	}

	tracelog.COMPLETED(service.UserId, "AddComment")
	return nil
}

func GetCommentByID(service *services.Service, id string) (comment *models.Comment, err error) {
	defer helper.CatchPanic(&err, service.UserId, "GetCommentByID")

	tracelog.STARTED(service.UserId, "GetCommentByID")

	comment = &models.Comment{}
	err = service.DBAction("comments",
		func(collection *mgo.Collection) error {
			return collection.FindId(bson.ObjectIdHex(id)).One(comment)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetArticleByID")
		return nil, fmt.Errorf("Unable find comment by ID: %s \nError: %s", id, err)
	}

	tracelog.COMPLETED(service.UserId, "GetCommentByID")
	return comment, nil
}

func AddChildComment(service *services.Service, parentId string, comment models.Comment) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "AddChildComment")

	var parent *models.Comment
	parent, err = GetCommentByID(service, parentId)
	comment.Path = strings.Join([]string{parent.Path, parentId}, ":")

	err = AddComment(service, comment)

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetUserArticles")
		return fmt.Errorf("Unable add child comment: %s \nError: %s", comment.ID, err)
	}

	return nil
}

func DeleteComment(service *services.Service, id string) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "DeleteComment")

	err = service.DBAction("comments", func(collection *mgo.Collection) error {
		return collection.RemoveId(bson.ObjectIdHex(id))
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "DeleteComment")
		return fmt.Errorf("Unable delete comment by id: %s \nError: %s", id, err)
	}

	return nil
}

// func UpdateText(service *services.Service, id, text string) (err error) {
// 	defer helper.CatchPanic(&err, service.UserId, "UpdateText")

// 	um := bson.M{"text": text}
// 	err = service.DBAction("articles", func(collection *mgo.Collection) error {
// 		return collection.UpdateId(bson.ObjectIdHex(id), um)
// 	})

// 	if err != nil {
// 		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpdateText")
// 		return fmt.Errorf("Unable update text of article by id: %s \nError: %s", id, err)
// 	}

// 	return nil
// }

// func UpdateTitle(service *services.Service, id, title string) (err error) {
// 	defer helper.CatchPanic(&err, service.UserId, "UpdateText")

// 	um := bson.M{"title": title}
// 	err = service.DBAction("articles", func(collection *mgo.Collection) error {
// 		return collection.UpdateId(bson.ObjectIdHex(id), um)
// 	})

// 	if err != nil {
// 		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpdateText")
// 		return fmt.Errorf("Unable update text of article by id: %s \nError: %s", id, err)
// 	}

// 	return nil
// }

// func UpdateArticle(service *services.Service, article models.Article) (err error) {
// 	defer helper.CatchPanic(&err, service.UserId, "UpdateArticle")

// 	err = service.DBAction("articles", func(collection *mgo.Collection) error {
// 		return collection.UpdateId(article.ID, article)
// 	})

// 	if err != nil {
// 		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpdateText")
// 		return fmt.Errorf("Unable update article ID: %s \nError: %s", article.ID, err)
// 	}

// 	return nil
// }

// func ChangeRating(service *services.Service, articleId string, value int) (err error) {
// 	defer helper.CatchPanic(&err, service.UserId, "ChangeRating")

// 	var article *models.Article

// 	article, err = GetArticleByID(service, articleId)
// 	rating := &article.Rating
// 	rating.ChancgeRating(value)

// 	err = UpdateArticle(service, *article)

// 	if err != nil {
// 		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "ChangeRating")
// 		return fmt.Errorf("Unable update article ID: %s \nError: %s", articleId, err)
// 	}

// 	return nil
// }
