package helper

import (
	// "fmt"
	// "github.com/rpoletaev/rev-dirty-chat/app/models"
	"errors"
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
		case "time":
			val, err := time.Parse("01 January, 2006", changeValue)
			if err != nil {
				return nil, err
			}

			return bson.M{colName: val}, nil
		default:
			return nil, errors.New("Неподдерживаемый тип данных")
		}
	} else {
		return nil, errors.New("Не удалось найти поле")
	}
}

// func ParseValue(val interface{}) interface{} {

// }
