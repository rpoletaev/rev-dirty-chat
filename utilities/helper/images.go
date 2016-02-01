package helper

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"image"
	"io/ioutil"
	"runtime"
	"strings"
)

func ProcessAvatar(img image, width int, height int, dir, fname string) path string {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	imgNum := GetImageNumber(dir, fname)
	resized := imaging.Resize(img, width, height, imaging.Gaussian)
	target := fmt.Sprintf("%s/%s.png", dir, imgNum)
	imaging.Save(resized, target)
}

func CreateAvatar(img image, fname string) (small, big string) {
	dir := fmt.Sprintf("%s/img/%s/avatar", PUBLIC_PATH, fname)
	small = ProcessAvatar(img, 0, 80, dir, fname)
	tracelog.TRACE(MAIN_GO_ROUTINE, "CreateAvatar", "[%s]", small)
	
	dir := fmt.Sprintf("%s/img/%s", PUBLIC_PATH, fname)
	big = ProcessAvatar(img, 300, 0, dir, fname)
	tracelog.TRACE(MAIN_GO_ROUTINE, "CreateAvatar", "[%s]", big)
	return small, big
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