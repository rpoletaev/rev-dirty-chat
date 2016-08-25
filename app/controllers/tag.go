package controllers

import (
	"strings"

	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/articles"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
)

type Tag struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*Tag).Before, revel.BEFORE)
	revel.InterceptMethod((*Tag).After, revel.AFTER)
	revel.InterceptMethod((*Tag).Panic, revel.PANIC)
}

func (tag *Tag) Index() revel.Result {
	tags, err := articles.GetAllTags(tag.Services())
	if err != nil {
		return tag.RenderTemplate("errors/404.html")
	}

	return tag.Render(tags)
}

func (tag *Tag) Create(id, synonim string) revel.Result {
	println("synonims is ", synonim)
	synonims := strings.Split(synonim, ",")
	for _, syn := range synonims {
		println(strings.TrimSpace(syn))
	}

	synonims = append(synonims, id)
	allSynonims, err := articles.GetAllSynonims(tag.Services(), synonims)
	if err != nil {
		println(err)
		return tag.RenderTemplate("errors/404.html")
	}

	for _, srcSyn := range synonims {
		if !helper.StringsContains(allSynonims, srcSyn, true) {
			allSynonims = append(allSynonims, strings.ToLower(srcSyn))
		}
	}

	for _, syn := range allSynonims {
		newTag := models.Tag{ID: syn, Synonims: helper.StringsWithoutFirstEntry(allSynonims, syn)}
		articles.InsertTag(tag.Services(), newTag)
	}

	return tag.RenderText("tag synonims is ", allSynonims)
}
