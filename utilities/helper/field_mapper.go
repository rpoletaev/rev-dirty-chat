package helper

import (
	"errors"
	"fmt"
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strconv"
	"time"
)

func GetChangesMap(targetStruct interface{}, changeName, changeValue string) (map[string]interface{}, error) {
	structType := reflect.TypeOf(targetStruct)
	field, found := structType.FieldByName(changeName)

	if found {
		colName := string(field.Tag.Get("bson"))
		if colName == "" {
			panic("The Field has'nt db mapping")
		}

		switch field.Type.Name() {
		case "bool":
			val, err := strconv.ParseBool(changeValue)
			if err != nil {
				return nil, err
			}

			return bson.M{colName: val}, nil
		case "string":
			return bson.M{colName: changeValue}, nil
		case "int", "int64":
			val, err := strconv.ParseInt(changeValue, 10, 64)
			if err != nil {
				return nil, err
			}

			return bson.M{colName: val}, nil
		case "float", "float64":
			val, err := strconv.ParseFloat(changeValue, 64)
			if err != nil {
				return nil, err
			}

			return bson.M{colName: val}, nil
		case "Time":
			val, err := time.Parse("02 January, 2006", changeValue)
			if err != nil {
				return nil, err
			}

			return bson.M{colName: val}, nil
		case "Sex":
			sexes := models.GetSexes()
			if _, ok := sexes[changeValue]; !ok {
				return nil, errors.New("Пол не найден")
			}

			return bson.M{colName: sexes[changeValue]}, nil
		case "Position":
			positions := models.GetPositions()
			if _, ok := positions[changeValue]; !ok {
				return nil, errors.New("Позиционирование не найдено")
			}

			return bson.M{colName: positions[changeValue]}, nil
		case "Orientation":
			orientations := models.GetOrientations()
			if _, ok := orientations[changeValue]; !ok {
				return nil, errors.New("Ориентация не найдена")
			}

			return bson.M{colName: orientations[changeValue]}, nil
		default:
			return nil, errors.New(fmt.Sprintln(field.Type.Name(), "- неподдерживаемый тип данных"))
		}
	} else {
		return nil, errors.New("Не удалось найти поле")
	}
}

// func ParseValue(val interface{}) interface{} {

// }
