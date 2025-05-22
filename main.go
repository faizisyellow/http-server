package main

import (
	"github.com/faizisyellow/http-server/http"
)

func main() {

	http.HandleFunc("/", base)
	// http.HandleFunc("/echo/{variable}", echo)
	// http.HandleFunc("/user-agent", userAgent)
	// http.HandleFunc("/files/{file}", fileOperation)

	http.ListenAndServe("4221")

}
