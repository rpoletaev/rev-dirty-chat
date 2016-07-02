package controllers

import (
	"time"

	"strings"

	"fmt"
	"html/template"

	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/articles"
	"github.com/russross/blackfriday"
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*Article).Before, revel.BEFORE)
	revel.InterceptMethod((*Article).After, revel.AFTER)
	revel.InterceptMethod((*Article).Panic, revel.PANIC)
}

func (art *Article) New() revel.Result {
	if !art.Authenticated() {
		return art.Redirect("/session/new")
	}

	return art.Render()
}

func (art *Article) Edit(id string) revel.Result {
	article, err := articles.GetArticleByID(art.Services(), id)
	if err != nil {
		return art.RenderTemplate("errors/404.html")
	}

	if article.Author.ID != art.Session["CurrentUserID"] {
		return art.RenderTemplate("errors/500.html")
	}

	return art.Render(article)
}

func (art *Article) Create(title, source_text, article_source, tags string) revel.Result {
	if !art.Authenticated() {
		return art.Redirect("/session/new")
	}

	fmt.Println(article_source)
	tgs := strings.Split(tags, "#")
	tagIds := []string{}
	for _, tag := range tgs {
		if strings.TrimSpace(tag) != "" {
			tagIds = append(tagIds, tag)
		}
	}

	fmt.Println("sourse text is", source_text)
	text := template.HTML(blackfriday.MarkdownCommon([]byte(source_text))[:])
	fmt.Println(text)
	article := models.Article{
		ID:    bson.NewObjectId(),
		Title: title,
		Tags:  tagIds,
		Post: models.Post{
			Author: models.ShortUser{
				ID:    art.Session["CurrentUserID"],
				Login: art.Session["Login"],
			},
			Rating:    models.Rating{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Markdown:  source_text,
			Text:      text,
		},
	}

	err := articles.InsertArticle(art.Services(), article)
	if err != nil {
		return art.RenderTemplate("errors/404.html")
	}

	return art.Redirect(fmt.Sprintf("/articles/show/%s", article.ID.Hex()))
}

func (art *Article) Show(id string) revel.Result {
	article, err := articles.GetArticleByID(art.Services(), id)
	fmt.Println(article)
	//article.Text = html.scapeString(article.Text)
	//fmt.Println(article.Text)
	if err != nil {
		return art.RenderTemplate("errors/404.html")
	}

	return art.Render(article)
}

func (art *Article) Index(filter map[string]interface{}) revel.Result {
	articles, count, err := articles.GetPageArticles(art.Services(), filter, 1, 10)

	if err != nil {
		return art.RenderError(err)
	}

	pageCount := count % 10
	return art.Render(articles, pageCount)
}

// func (art *Article) Delete() revel.Result {

// }

// func (art *Article) Update() revel.Result {

// }
