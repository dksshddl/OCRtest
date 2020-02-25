package utils

import (
	"encoding/json"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type OCRRequest struct {
	Version   string     `json:"version"`
	RequestId string     `json:"requestId"`
	Timestamp string     `json:"timestamp"`
	Lang      string     `json:"lang, omitempty"`
	Images    []OCRImage `json:"images"`
}

type OCRImage struct {
	Format     string   `json:"format"`
	Data       string   `json:"data,omitempty"`
	Url        string   `json:"url,omitempty"`
	Name       string   `json:"name"`
	TemplateId []string `json:"templateId, omitempty"`
}

func ConvertImage(f io.Reader, suffix string, o string) *image.NRGBA {
	var img image.Image
	var err error

	switch suffix {
	case "jpg":
		img, err = jpeg.Decode(f)
		println("convert file to jpg image")
		if err != nil {
			return nil
		}
	case "png":
		img, err = png.Decode(f)
		println("convert file to png image")
		if err != nil {
			return nil
		}
	case "gif":
		img, err = gif.Decode(f)
		println("convert file to gif image")
		if err != nil {
			return nil
		}
	}
	switch o {
	case "1":
		return imaging.Clone(img)
	case "2":
		return imaging.FlipV(img)
	case "3":
		return imaging.Rotate180(img)
	case "4":
		return imaging.Rotate180(imaging.FlipV(img))
	case "5":
		return imaging.Rotate270(imaging.FlipV(img))
	case "6":
		return imaging.Rotate270(img)
	case "7":
		return imaging.Rotate90(imaging.FlipV(img))
	case "8":
		return imaging.Rotate90(img)
	}
	return imaging.Clone(img)
}

func RequsetOCR(data string) {
	ocrURl := "https://f150jn75jw.apigw.ntruss.com/custom/v1/628/2b6ca66f008b370fad3002b8c5d1b5a89c90fcd9fb5f5e7ed909df1e80d13b5e/general"
	ocrSecretKey := "YXhlaERrVnFpcFNPTENRbWtCUkJHSVRJTkJ4V21CYmM="
	timestamp := int(time.Now().Unix())
	ocrImages := make([]OCRImage, 1)
	ocrImages[0] = OCRImage{
		Format: "jpg",
		Data:   data,
		Name:   "sample",
	}
	ocrRequest := OCRRequest{
		Version:   "V1",
		RequestId: uuid.New().String(),
		Timestamp: strconv.Itoa(timestamp),
		Lang:      "ko",
		Images:    ocrImages,
	}
	doc, _ := json.Marshal(ocrRequest)
	println(string(doc))
	req, _ := http.NewRequest("POST", ocrURl, strings.NewReader(string(doc)))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-OCR-SECRET", ocrSecretKey)

	client := &http.Client{}
	res, _ := client.Do(req)
	resBody, _ := ioutil.ReadAll(res.Body)
	println(string(resBody))
}
