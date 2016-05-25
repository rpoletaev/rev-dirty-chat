package controllers

import (
	//"container/list"
	"fmt"

	"github.com/revel/revel"
	cb "github.com/rpoletaev/rev-dirty-chat/app/controllers/base"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services/regionService"
	"github.com/rpoletaev/rev-dirty-chat/utilities/regions"
	//"gopkg.in/mgo.v2/bson"
)

type Region struct {
	cb.BaseController
}

func init() {
	revel.InterceptMethod((*Region).Before, revel.BEFORE)
	revel.InterceptMethod((*Region).After, revel.AFTER)
	revel.InterceptMethod((*Region).Panic, revel.PANIC)
}

func (c *Region) Load() revel.Result {
	if !c.IsAdmin() {
		return c.RenderTemplate("templates/errors/500.html")
	}

	return c.RenderText("Регионы загружены")
	region_file := fmt.Sprintf("%s/regions.xlsx", revel.BasePath)
	fmt.Println(region_file)
	fmt.Println(revel.BasePath)
	lst, err := regions.GetRegionsFromFile(region_file)
	if err != nil {
		return c.RenderError(err)
	}

	for e := lst.Front(); e != nil; e = e.Next() {
		region := models.Region{
			Name: e.Value.(string),
		}
		err = regionService.InsertRegion(c.Services(), &region)
		if err != nil {
			return c.RenderError(err)
		}
	}
	return c.RenderText("Все заебись")
}

func (c *Region) GetRegions() revel.Result {
	rgns, err := regionService.GetAllRegions(c.Services())
	if err != nil {
		return c.RenderTemplate("templates/Errors/404.html")
	}

	return c.RenderJson(rgns)
}

// func (c *Region) CreateRooms() revel.Result {
// 	err := regionService.CreateRegionRooms(c.Services())
// 	if err != nil {
// 		c.RenderError(err)
// 	}
// 	return c.RenderText("Все заебись")
// }
