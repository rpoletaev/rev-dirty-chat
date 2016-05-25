package regionService

import (
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

func InsertRegion(service *services.Service, region *models.Region) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "Insert Region")

	err = service.DBAction("regions",
		func(collection *mgo.Collection) error {
			return collection.Insert(region)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "Insert Region")
		return err
	}

	tracelog.COMPLETED(service.UserId, "Insert Region")
	return nil
}

func GetAllRegions(service *services.Service) (regions *[]models.Region, err error) {
	defer helper.CatchPanic(&err, service.UserId, "GetRegions")
	regions = &[]models.Region{}
	err = service.DBAction("regions",
		func(collection *mgo.Collection) error {
			return collection.Find(nil).All(regions)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "GetRegions")
		return nil, err
	}

	tracelog.COMPLETED(service.UserId, "GetRegions")
	return regions, nil
}

// func CreateRegionRooms(service *services.Service) (err error) {
// 	return nil
// 	var regions *[]models.Region
// 	regions, err = GetAllRegions(service)
// 	if err != nil {
// 		return err
// 	}

// 	for _, region := range *regions {
// 		header := models.RoomHeader{
// 			Name:      region.Name,
// 			IsPrivate: false,
// 			Region:    region.ID.Hex(),
// 		}

// 		chatService.InsertRoom(service, &room)
// 	}
// 	return nil
// }
