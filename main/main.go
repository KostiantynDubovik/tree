package main

import (
	"net/http"
	"html/template"

	"image"
	"image/jpeg"
	_"image/png"
	_"image/gif"
)

type page struct {
	Title string
	Msg   string
}

func index(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-type", "text/html")

	title := request.URL.Path[len("/"):]

	if title != "exec/" {
		fileTemplate, err := template.ParseFiles("src/web/html/index.html")
		if err != nil {
			defer errorHandler(responseWriter, "error in sending html")
		}
		fileTemplate.Execute(responseWriter, &page{Title: "Convert Image"})

	} else {
		imageFile, fileHeader, err := request.FormFile("imageFile")
		if err != nil {
			defer errorHandler(responseWriter,"error in image recieving")
		}
		decodedImage, fileName, err := image.Decode(imageFile)
		if err != nil {
			defer errorHandler(responseWriter, "error in image decoding")
		}
		responseWriter.Header().Set("Content-type", "image/jpeg")
		responseWriter.Header().Set("Content-Disposition", "filename=\""+fileHeader.Filename+"."+fileName+"\"")
		err = jpeg.Encode(responseWriter, decodedImage, &jpeg.Options{0})
		if err!=nil{
			defer errorHandler(responseWriter,"error in file format")
		}

	}
}

func errorHandler(writer http.ResponseWriter, blockName string) {
	recover()
	writer.Header().Set("Content-type", "text/html")
	writer.Write([]byte(blockName))
}

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(":80", nil)
}
