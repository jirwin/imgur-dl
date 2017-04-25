package imgur

import (
	"net/http"
	"time"

	"encoding/json"
	"fmt"

	"io"
	"os"

	"github.com/parnurzeal/gorequest"
)

type Imgur struct {
	clientId   string
	httpClient *http.Client
}

func (i *Imgur) makeRequest(endpoint string, obj interface{}) error {
	request := gorequest.New()
	_, body, errs := request.Get(fmt.Sprintf("https://api.imgur.com/3/%s", endpoint)).
		Set("Authorization", fmt.Sprintf("Client-ID %s", i.clientId)).
		End()
	if len(errs) != 0 {
		return errs[0]
	}

	resp := &ImgurResponse{
		Data: obj,
	}

	err := json.Unmarshal([]byte(body), resp)
	if err != nil {
		return err
	}

	return nil
}

func (i *Imgur) GetAlbum(albumId string) (*Album, error) {
	album := &Album{}

	err := i.makeRequest(fmt.Sprintf("album/%s", albumId), album)
	if err != nil {
		return nil, err
	}

	return album, nil
}

func (i *Imgur) GetGallery(galleryId string) (*Gallery, error) {
	gallery := &Gallery{}

	err := i.makeRequest(fmt.Sprintf("gallery/album/%s", galleryId), gallery)
	if err != nil {
		return nil, err
	}

	return gallery, nil
}

func (i *Imgur) DownloadImage(url, writePath string) error {
	if _, err := os.Stat(writePath); err == nil {
		return nil
	}

	response, err := i.httpClient.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	file, err := os.Create(writePath + ".incomplete")
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	err = os.Rename(writePath+".incomplete", writePath)
	if err != nil {
		return err
	}

	return nil
}

func MakeImgur(clientId string) *Imgur {
	httpClient := &http.Client{
		Timeout: time.Second * 100,
	}

	return &Imgur{
		clientId:   clientId,
		httpClient: httpClient,
	}
}
