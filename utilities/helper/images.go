package helper

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"io/ioutil"
	"runtime"
	"strings"
)

func ProcessAvatar(img image.Image, width int, height int, dir, fname string) (path string) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	imgNum := GetImageNumber(dir)
	resized := imaging.Resize(img, width, height, imaging.Gaussian)
	path = fmt.Sprintf("%s/%d.png", dir, imgNum)
	imaging.Save(resized, path)
	return path
}

func CreateAvatar(img image.Image, fname, base_path string) (small, big string) {
	dir := fmt.Sprintf("%s/public/img/%s/avatar", base_path, fname)
	small = ProcessAvatar(img, 0, 80, dir, fname)
	small = strings.Replace(small, base_path, "", 1)

	dir = fmt.Sprintf("%s/public/img/%s", base_path, fname)
	big = ProcessAvatar(img, 300, 0, dir, fname)
	big = strings.Replace(big, base_path, "", 1)
	return small, big
}

func GetImageNumber(dir string) int {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	return len(files)
}
