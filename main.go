package main

import (
	"github.com/faizisyellow/http-server/http"
)

func main() {

	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/echo/{text}", echoHandler)
	http.HandleFunc("/user-agent", userAgentHandler)
	http.HandleFunc("/files/*file", fileHandler)

	http.ListenAndServe("4221")

}
