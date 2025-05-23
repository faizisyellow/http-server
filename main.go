package main

import (
	"github.com/faizisyellow/http-server/http"
)

func main() {

	http.HandleFunc("/", baseHandler, "GET")
	http.HandleFunc("/echo/{text}", echoHandler, "GET")
	http.HandleFunc("/user-agent", userAgentHandler, "GET")
	http.HandleServer("/files", fileHandler, "GET")
	http.HandleServer("/files", fileUploadHandler, "POST")

	http.ListenAndServe("4221")

}
