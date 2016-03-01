package helper

import (
	"fmt"
	"os"
)

func CreateUserFS(base_path, login string) {
	img_dir := fmt.Sprintf("%s/public/img/%s/avatar", base_path, login)
	err := os.MkdirAll(img_dir, 0777)
	if err != nil {
		panic(err.Error())
	}
}
