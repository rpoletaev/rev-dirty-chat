package articles

import (
	"fmt"

	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InsertArticle(service *services.Service, article models.Article) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "InsertArticle")

	tracelog.STARTED(service.UserId, "InsertArticle")

	// tracelog.TRACE(helper.MAIN_GO_ROUTINE, "InsertArticle", "Query : %s", mongo.ToString(queryMap))
	// Execute the query
	err = service.DBAction("articles",
		func(collection *mgo.Collection) error {
			return collection.Insert(article)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "InsertArticle")
		return err
	}

	tracelog.COMPLETED(service.UserId, "InsertArticle")
	return nil
}

func GetArticleByID(service *services.Service, id string) (article *models.Article, err error) {
	defer helper.CatchPanic(&err, service.UserId, "GetArticleByID")

	tracelog.STARTED(service.UserId, "GetArticleByID")

	article = &models.Article{}
	err = service.DBAction("articles",
		func(collection *mgo.Collection) error {
			return collection.FindId(bson.ObjectIdHex(id)).One(article)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetArticleByID")
		return nil, fmt.Errorf("Unable find article by ID: %s \nError: %s", id, err)
	}

	tracelog.COMPLETED(service.UserId, "GetArticleByID")
	return article, nil
}

func GetUserArticles(service *services.Service, userId string) (articles *[]models.Article, err error) {
	defer helper.CatchPanic(&err, service.UserId, "GetUserArticles")

	qm := bson.M{"post.author.id": userId}
	articles = &[]models.Article{}
	err = service.DBAction("articles",
		func(collection *mgo.Collection) error {
			return collection.Find(qm).All(articles)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetUserArticles")
		return nil, fmt.Errorf("Unable find articles of user: %s \nError: %s", userId, err)
	}

	return articles, nil
}

func DeleteArticle(service *services.Service, id string) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "DeleteArticle")

	err = service.DBAction("articles", func(collection *mgo.Collection) error {
		return collection.RemoveId(bson.ObjectIdHex(id))
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "DeleteArticle")
		return fmt.Errorf("Unable delete article by id: %s \nError: %s", id, err)
	}

	return nil
}

func UpdateText(service *services.Service, id, text string) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "UpdateText")

	um := bson.M{"text": text}
	err = service.DBAction("articles", func(collection *mgo.Collection) error {
		return collection.UpdateId(bson.ObjectIdHex(id), um)
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpdateText")
		return fmt.Errorf("Unable update text of article by id: %s \nError: %s", id, err)
	}

	return nil
}

func UpdateTitle(service *services.Service, id, title string) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "UpdateText")

	um := bson.M{"title": title}
	err = service.DBAction("articles", func(collection *mgo.Collection) error {
		return collection.UpdateId(bson.ObjectIdHex(id), um)
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpdateText")
		return fmt.Errorf("Unable update text of article by id: %s \nError: %s", id, err)
	}

	return nil
}

func UpdateArticle(service *services.Service, article models.Article) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "UpdateArticle")

	err = service.DBAction("articles", func(collection *mgo.Collection) error {
		return collection.UpdateId(article.ID, article)
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "UpdateText")
		return fmt.Errorf("Unable update article ID: %s \nError: %s", article.ID, err)
	}

	return nil
}

func GetPageArticles(service *services.Service, filter bson.M, page, count int) (articles *[]models.Article, filteredCount int, err error) {
	defer helper.CatchPanic(&err, service.UserId, "GetArticles")

	skipingCount := (page - 1) * count
	articles = &[]models.Article{}

	err = service.DBAction("articles", func(collection *mgo.Collection) error {
		query := collection.Find(filter)
		filteredCount, err = query.Count()
		if err != nil {
			return err
		}

		return query.Sort("-post.updatedAt").Skip(skipingCount).Limit(count).All(articles)
	})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetArticles")
		return nil, 0, fmt.Errorf("Unable Find articles to page: %d \nError: %s", page, err)
	}

	return articles, filteredCount, nil
}

func ChangeRating(service *services.Service, articleId string, value int) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "ChangeRating")

	var article *models.Article

	article, err = GetArticleByID(service, articleId)
	rating := &article.Rating
	rating.ChangeRating(value)

	err = UpdateArticle(service, *article)

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "ChangeRating")
		return fmt.Errorf("Unable update article ID: %s \nError: %s", articleId, err)
	}

	return nil
}
