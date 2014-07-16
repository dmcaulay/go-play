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
	resp, err := http.Get("http://extcam-14.se.axis.com/axis-cgi/mjpg/video.cgi?resolution=640x480")
	if err != nil {
		log.Fatal(err)
	} else {
		mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(mediaType, "multipart/") {
			mr := multipart.NewReader(resp.Body, params["boundary"])
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
				fmt.Printf("MJPEG frame %d: %d\n", i, len(mjpeg))
			}
		}
	}
}
