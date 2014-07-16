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
	urls := [1]string{"http://extcam-14.se.axis.com/axis-cgi/mjpg/video.cgi?resolution=640x480"}
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
				readStream(id, multipart.NewReader(resp.Body, params["boundary"]))
			}
		}
	}
}

func readStream(id int, mr *multipart.Reader) {
	for i := 0; true; i++ {
		part, err := mr.NextPart()
		if err == io.EOF {
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
