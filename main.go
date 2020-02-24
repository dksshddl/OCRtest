package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type OCRRequest struct {
	Version   string    `json:"version"`
	RequestId string    `json:"requestId"`
	Timestamp string    `json:"timestamp"`
	Lang      string    `json:"lang"`
	Images    *OCRImage `json:"images"`
}

type OCRImage struct {
	Format     string   `json:"format"`
	Url        string   `json:"url"`
	Data       string   `json:"data"`
	Name       string   `json:"name"`
	TemplateId []string `json:"templateId"`
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html")
	t.Execute(w, "Hello world!")
}

func OCRtest(w http.ResponseWriter, r *http.Request) {
	ocrURl := "https://f150jn75jw.apigw.ntruss.com/custom/v1/628/2b6ca66f008b370fad3002b8c5d1b5a89c90fcd9fb5f5e7ed909df1e80d13b5e/general"
	ocrSecretKey := "YXhlaERrVnFpcFNPTENRbWtCUkJHSVRJTkJ4V21CYmM="
	_ = r.ParseForm()
	println(r)
	file, _, err := r.FormFile("sample")
	if err != nil {
		panic(err)
	} else {
		fmt.Println(file)
	}
	timestamp := int(time.Now().Unix())

	ocrImages := OCRImage{
		Format:     "jpg",
		Url:        "localhost:8080/sample.jpg",
		Data:       "",
		Name:       "sample",
		TemplateId: []string{"test"},
	}
	ocrRequset := OCRRequest{
		Version:   "V1",
		RequestId: uuid.New().String(),
		Timestamp: strconv.Itoa(timestamp),
		Lang:      "ko",
		Images:    &ocrImages,
	}
	doc, err := json.Marshal(ocrRequset)
	println(doc)
	req, err := http.NewRequest("POST", ocrURl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-OCR-SECRET", ocrSecretKey)

}

func main() {
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	http.HandleFunc("/index", index)
	http.HandleFunc("/sample", OCRtest)
	server.ListenAndServe()
}
