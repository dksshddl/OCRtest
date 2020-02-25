package main

import (
	"OCRtest/utils"
	"encoding/base64"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		OCRtest(w, r)
	case "GET":
		t, _ := template.ParseFiles("template/index.html")
		_ = t.Execute(w, "Upload Image")
	}
}

func OCRtest(w http.ResponseWriter, r *http.Request) {
	ch := make(chan string)
	var orientation *tiff.Tag

	file, fileHeader, _ := r.FormFile("sample")
	defer file.Close()
	split := strings.Split(fileHeader.Filename, ".")
	suffix := split[len(split)-1]

	go func() {
		x, _ := exif.Decode(file)
		_, _ = file.Seek(0, 0)
		orientation, _ = x.Get(exif.Orientation)
		ch <- orientation.String()
	}()

	go func() {
		val := <-ch
		img := utils.ConvertImage(file, suffix, val)
		imaging.Save(img, "./images/"+fileHeader.Filename)
		bytesData, _ := ioutil.ReadFile("./images/" + fileHeader.Filename)
		encData := base64.StdEncoding.EncodeToString(bytesData)
		utils.RequsetOCR(encData)
	}()
}

func main() {
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	http.HandleFunc("/index", index)
	server.ListenAndServe()
}
