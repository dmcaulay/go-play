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

type frame struct {
	id   int
	num  int
	jpeg []byte
}

func main() {
	urls := []string{
		"http://extcam-14.se.axis.com/axis-cgi/mjpg/video.cgi?resolution=320x240",
		"http://extcam-12.se.axis.com/axis-cgi/mjpg/video.cgi?resolution=320x240",
	}

	c := make(chan frame)
	for id, url := range urls {
		// use mjpeg lib
		go startStream(id, url, c)
	}

	// handle frames
	for {
		f := <-c
		fileName := fmt.Sprintf("frame.%d.%d.jpg", f.id, f.num)
		log.Println(fileName)
		ioutil.WriteFile(fileName, f.jpeg, 0644)
	}
}

// mjpeg lib
func startStream(id int, url string, c chan<- frame) {
	log.Printf("starting stream %d", id)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		readStream(id, multipart.NewReader(resp.Body, params["boundary"]), c)
	}
}

func readStream(id int, mr *multipart.Reader, c chan<- frame) {
	for i := 0; true; i++ {
		part, err := mr.NextPart()
		if err == io.EOF {
			log.Printf("ending stream %d", id)
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		jpeg, err := ioutil.ReadAll(part)
		if err != nil {
			log.Fatal(err)
		}
		c <- frame{id: id, num: i, jpeg: jpeg}
	}
}
