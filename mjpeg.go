package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

func main() {
	urls := []string{
		"http://extcam-14.se.axis.com/axis-cgi/mjpg/video.cgi?resolution=320x240",
		"http://extcam-12.se.axis.com/axis-cgi/mjpg/video.cgi?resolution=320x240",
	}
	for id, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		} else {
			mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
			if err != nil {
				log.Fatal(err)
			}
			if strings.HasPrefix(mediaType, "multipart/") {
				log.Printf("starting stream %d", id)
				go readStream(id, multipart.NewReader(resp.Body, params["boundary"]))
			}
		}
	}
	var end string
	fmt.Scanln(&end)
}

func readStream(id int, mr *multipart.Reader) {
	for i := 0; true; i++ {
		part, err := mr.NextPart()
		if err == io.EOF {
			log.Printf("ending stream %d", id)
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		mjpeg, err := ioutil.ReadAll(part)
		if err != nil {
			log.Fatal(err)
		}
		fileName := fmt.Sprintf("frame.%d.%d.jpg", id, i)
		log.Println(fileName)
		ioutil.WriteFile(fileName, mjpeg, 0644)
	}
}
