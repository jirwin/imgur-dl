package main

import (
	"fmt"

	pb "gopkg.in/cheggaaa/pb.v1"

	"os"

	"path/filepath"

	"errors"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jirwin/imgur-dl/imgur"
	"github.com/urfave/cli"
)

const Version = "0.0.3"

func run(c *cli.Context) error {
	var clientId string
	if !c.IsSet("clientId") {
		return errors.New("Missing required argument: --clientId")
	} else {
		clientId = c.String("clientId")
	}

	var url string
	if !c.IsSet("url") {
		return errors.New("Missing required argument: --url")
	} else {
		url = c.String("url")
	}

	var skipNsfw bool
	if c.IsSet("skip-nsfw") {
		skipNsfw = c.Bool("skip-nsfw")
	}

	client := imgur.MakeImgur(clientId)

	urlSplit := strings.Split(url, "/")
	if len(urlSplit) == 0 {
		return fmt.Errorf("Unable to parse gallery url: %s", url)
	}

	galleryId := urlSplit[len(urlSplit)-1]

	gallery, err := client.GetGallery(galleryId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	downloadPath := filepath.Join(".", gallery.Id)
	os.MkdirAll(downloadPath, os.ModePerm)

	concurrency := 35
	sem := make(chan bool, concurrency)

	bar := pb.StartNew(gallery.ImagesCount)

	for _, img := range gallery.Images {
		sem <- true
		go func(img *imgur.Image) {
			defer func() {
				<-sem
				bar.Increment()
			}()

			if skipNsfw && img.Nsfw {
				log.Infof("Skipping nsfw image: %s", img.Link)
				return
			}

			err := client.DownloadImage(img.Link, filepath.Join(downloadPath, img.Id+".jpg"))
			if err != nil {
				log.Errorf("Unable to download image: %s", img.Link)
			}

		}(img)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	bar.FinishPrint("Gallery successfully downloaded")

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "imgur-dl"
	app.Version = Version
	app.Usage = "Download all of the images in a imgur gallery."
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Usage: "The url of the gallery to download.",
		},
		cli.BoolFlag{
			Name:  "skip-nsfw",
			Usage: "Skip images that are flagged as NSFW",
		},
		cli.StringFlag{
			Name:  "clientId",
			Usage: "The imgur app client id to use.",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
