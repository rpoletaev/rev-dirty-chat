package helper

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"io/ioutil"
	"runtime"
)

func CreateAvatar() (original, small string) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ioutil.
	img, err := imaging.Open("/home/roman/go/src/github.com/rpoletaev/dirty-chat/assets/img/noavatar.png")
	if err != nil {
		fmt.Println(err.Error())
	}

	resized := imaging.Resize(img, 0, 80, imaging.Gaussian)
	imaging.Save(resized, "/home/roman/go/src/github.com/rpoletaev/dirty-chat/assets/img/noavatar_resized_gausian.png")
	// blur_png := imaging.Blur(resized, 0.3)
	// imaging.Save(blur_png, "/home/roman/go/src/github.com/rpoletaev/dirty-chat/assets/img/rp_blur.jpg")
}

func RandomName(dir_name string) {

}
