package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	client_id     = "4c2bcf743596439"
	client_secret = "c01533b057713d7b592a64d915fab54b5b036782"
)

type ImgurResponse struct {
	Data    DataObject
	Success bool
	Status  int
}

type DataObject struct {
	Id             string
	Title          string
	Description    string
	DateTime       string
	Type           string
	Animated       bool
	Width          int
	Height         int
	Size           int
	Views          int
	Bandwidth      int
	Vote           bool
	Favorite       bool
	Nsfw           bool
	Section        bool
	AccountURL     string
	AccountID      int
	CommentPreview bool
	DeleteHash     string
	Name           string
	Link           string
}

func ImgurImageUpload(img []byte, description string) string {
	url := "https://api.imgur.com/3/image"
	imgString := base64.StdEncoding.EncodeToString(img)
	fmt.Println(imgString)
	reqJson := []byte(fmt.Sprintf(`{"image":"%s", "description":"%s"}`, imgString, description))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", "4c2bcf743596439"))
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response ImgurResponse
	json.Unmarshal(body, &response)

	if response.Success {
		return fmt.Sprintf("http://i.imgur.com/%s.png", response.Data.Id)
	} else {
		panic(fmt.Sprintf("unable to load image %s", response))
	}
}
