package http

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Header map[string]string

type Request struct {
	Method string
	Url    string
	Proto  string
	Header Header
	Body   []byte
}

type ResponseWrite struct {
	Conn net.Conn
}

func (r *ResponseWrite) Write(statusCode int, statusText string) {

	r.Conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %v %v\r\n", statusCode, statusText)))

}

type HandlerFunc func(Request, *ResponseWrite)

var routes = make(map[string]HandlerFunc)

func HandleFunc(pattern string, handler HandlerFunc) {
	routes[pattern] = handler
}

func ListenAndServe(p string) {

	port := fmt.Sprintf("0.0.0.0:%v", p)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer ln.Close()
	fmt.Printf("Listening on port :%v\n", p)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {

	data := make([]byte, 1024)

	n, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	req := parseHttpRequest(data[:n])

	fmt.Println(req)

	for pattern, handler := range routes {
		if pattern == req.Url {
			handler(req, &ResponseWrite{Conn: conn})
		}
	}

}

// TODO: refactor
func parseHttpRequest(data []byte) Request {
	reqs := strings.Split(string(data), "\r\n")
	body := strings.Split(string(data), "\r\n\r\n")[1]

	// make alocates the variable to memory
	var headers = make(Header)

	for _, v := range reqs {
		if strings.Contains(v, ":") {
			h := strings.Split(v, ":")
			if len(h) > 2 {
				headers[h[0]] = fmt.Sprintf("%v:%v", h[1], h[2])
			} else {
				headers[h[0]] = h[1]
			}
		}
	}

	return Request{
		Method: strings.Split(reqs[0], " ")[0],
		Url:    strings.Split(reqs[0], " ")[1],
		Proto:  strings.Split(reqs[0], " ")[2],
		Header: headers,
		Body:   []byte(body),
	}
}
