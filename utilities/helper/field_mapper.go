package helper

import (
	// "fmt"
	// "github.com/rpoletaev/rev-dirty-chat/app/models"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strconv"
	"time"
)

func GetChangesExpression(targetStruct interface{}, changeName, changeValue string) map[string]interface{} {
	structType := reflect.TypeOf(targetStruct)
	field, found := structType.FieldByName(changeName)

	if found {
		colName := string(field.Tag.Get("bson"))
		if colName == "" {
			panic("The Field has'nt db mapping")
		}

		switch field.Type.Name() {
		case "bool":
			val, err := strconv.ParseBool(str)
			if err != nil {
				panic("Переданное значение не соответствует типу данных")
			}

			return bson.M{colName: val}

		case "int" || "int64":
			val, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic("Переданное значение не соответствует типу данных")
			}

			return bson.M{colName: val}

		case "float" || "float64":
			val, err := strconv.ParseFloat(s, bitSize)
			if err != nil {
				panic("Переданное значение не соответствует типу данных")
			}

			return bson.M{colName: val}
		case "time":
			val, err := time.Parse("", val)
		}

	} else {
		panic("type unknown")
	}
}

// func ParseValue(val interface{}) interface{} {

// }
