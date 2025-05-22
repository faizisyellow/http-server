package main

import (
	"log"
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

	// TODO: figure it out the response
	var fileName string

	// 	// if the directory empty it return dot
	if filepath.Dir(dirpath) != "." {
		paths, err := filepath.Glob(dirpath + ".*")
		if err != nil {
			log.Fatal(err)
		}
		if len(paths) == 0 {
			w.Write(404, "Not Found", nil, "")
			return
		}

		fileName = paths[0]
	} else {

		w.Write(505, "Server Error", nil, "")
		return
	}
}
