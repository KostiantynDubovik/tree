package main

import (
	"net/http"
	"html/template"
	"image"
	"image/jpeg"
	_"image/png"
	_"image/gif"
	"github.com/gorilla/mux"
	"flag"
	"strconv"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

var id int = 0

var Nodes = make([]Node, 0)

type page struct {
	Title string
	Msg   string
}

func getTree() (interface{}, *handlerError) {
	return getTree(), nil
}

func addNode(responseWriter http.ResponseWriter, request *http.Request) (interface{}, *handlerError) {
	payload, e := parseNodeRequest(request)
	payload.NodeId = getNextId()
	if e != nil {
		return nil, e
	}

	imageFile, _, err := request.FormFile("imageFile")
	if err != nil {
		return nil, &handlerError{err, "Not a jpeg file", http.StatusUnsupportedMediaType}
	}
	SaveOriginalImage(imageFile, payload.ImageName)
	AddNode(payload.NodeId, payload.PrentId, payload.ImageName)

	return payload, nil
}
func removeNode(responseWriter http.ResponseWriter, request *http.Request) (interface{}, *handlerError) {
	param := mux.Vars(request)["nodeId"]
	nodeId, e := strconv.Atoi(param)
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}

	DeleteNode(nodeId)

	return make(map[string]string), nil
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
			defer errorHandler(responseWriter, "error in image recieving")
		}
		decodedImage, fileName, err := image.Decode(imageFile)
		if err != nil {
			defer errorHandler(responseWriter, "error in image decoding")
		}
		responseWriter.Header().Set("Content-type", "image/jpeg")
		responseWriter.Header().Set("Content-Disposition", "filename=\""+fileHeader.Filename+"."+fileName+"\"")
		err = jpeg.Encode(responseWriter, decodedImage, &jpeg.Options{Quality: 0})
		if err != nil {
			defer errorHandler(responseWriter, "error in file format")
		}

	}
}

func parseNodeRequest(request *http.Request) (Node, *handlerError) {
	data, e := ioutil.ReadAll(request.Body)
	if e != nil {
		return Node{}, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	var payload Node
	e = json.Unmarshal(data, &payload)
	if e != nil {
		return Node{}, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	return payload, nil
}

func errorHandler(writer http.ResponseWriter, blockName string) {
	recover()
	writer.Header().Set("Content-type", "text/html")
	writer.Write([]byte(blockName))
}

type handlerError struct {
	Error   error
	Message string
	Code    int
}

type handler func(responseWriter http.ResponseWriter, request *http.Request) (interface{}, *handlerError)

func (fn handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// here we could do some prep work before calling the handler if we wanted to

	// call the actual handler
	response, err := fn(responseWriter, request)

	// check for errors
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error)
		http.Error(responseWriter, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		log.Printf("ERROR: response from method is nil\n")
		http.Error(responseWriter, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	// turn the response into JSON
	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(responseWriter, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// send the response and log
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(bytes)
	log.Printf("%s %s %s %d", request.RemoteAddr, request.Method, request.URL, 200)
}

func main() {
	port := flag.Int("port", 80, "port to serve on")

	dir := flag.String("directory", "web/html/", "directory of web files")
	flag.Parse()

	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)

	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/static/", 302))
	router.Handle("/tree", handler(getTree)).Methods("GET")
	router.Handle("/tree", handler(addNode)).Methods("POST")
	router.Handle("/tree", handler(removeNode)).Methods("DELETE")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileHandler))
	http.Handle("/", router)
	address := fmt.Sprintf("127.0.0.1:%d", *port)

	//http.HandleFunc("/", index)
	RestyleImage("", "80x80")
	err := http.ListenAndServe(address, nil)
	fmt.Println(err.Error())
}

func getNextId() int {
	id += 1
	return id
}
