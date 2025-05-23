package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/faizisyellow/http-server/http"
)

func baseHandler(r http.Request, w http.ServerResponse) {

	w.Write(200, "OK", nil, "")
}

func echoHandler(r http.Request, w http.ServerResponse) {
	text := r.Params["text"]

	w.Write(200, "OK", nil, text)
}

func userAgentHandler(r http.Request, w http.ServerResponse) {

	w.Write(200, "OK", nil, r.Header["User-Agent"])
}

func fileHandler(r http.Request, w http.ServerResponse) {
	filename := r.FilePath
	if filename == "" {
		w.Write(400, "Bad Request", nil, "")
		return
	}

	fileDirectory := "files"
	fullPath := filepath.Join(fileDirectory, filename)

	_, err := os.Stat(fullPath)

	if err != nil {
		if os.IsNotExist(err) {
			w.Write(404, "Not Found", nil, "")
			return
		}
		w.Write(500, "Server Error", nil, "")
		return
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		w.Write(500, "Server Error", nil, "")
		return
	}

	w.Write(200, "OK", nil, string(content))
}

func fileUploadHandler(r http.Request, w http.ServerResponse) {
	filename := r.FilePath
	if filename == "" {
		w.Write(400, "Bad Request", nil, "")
		return
	}

	f, err := os.Create(fmt.Sprintf("files/%v", filename))
	if err != nil {
		w.Write(500, "Server Error", nil, "")
		return
	}

	defer f.Close()

	_, err = f.Write(r.Body)
	if err != nil {
		w.Write(500, "Server Error", nil, "")
		return
	}

	w.Write(201, "Created", nil, "")
}
