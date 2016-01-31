package helper

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"io/ioutil"
	"runtime"
	"strings"
)

func CreateAvatar(img image, fname string) (path string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	dir := fmt.Sprintf("%s/img/avatar", PUBLIC_PATH)
	imgNum := GetImageNumber(dir, fname)

	resized := imaging.Resize(img, 0, 80, imaging.Gaussian)
	imaging.Save(resized, fmt.Sprintf("%s/img/%s/avatar/%s.png", dir, fname, imgNum)
}

func CreateMainPhoto(img image, fname string)  path string{
	runtime.GOMAXPROCS(runtime.NumCPU())

	resized := imaging.Resize(img, 0, 80, imaging.Gaussian)
	imaging.Save(resized, "public/img/avatar/noavatar_resized_gausian.png")	
}

func GetImageNumber(dir, fname string) int{
	files, err := ioutil.ReadDir()
	if err != nil {
		panic(err)
	}

	counter := 0
	for _,file := range files {
		if strings.Contains(file.Name(), fname) {
			counter++
		}
	}

	return counter
}