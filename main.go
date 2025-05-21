package main

import (
	"fmt"
	"net"
	"strings"
)

type Header map[string]string

type Request struct {
	Method string
	Url    string
	Proto  string
	Header Header
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	go handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	data := make([]byte, 1024)

	_, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	req := ParseHttpRequest(data)

	if req.Url == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	} else if strings.HasPrefix(req.Url, "/echo/") {

		vars, _ := strings.CutPrefix(req.Url, "/echo/")
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(vars), vars)))

	} else if req.Url == "/user-agent" {
		userAgent := req.Header["User-Agent"]

		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(userAgent), userAgent)))
	} else {

		conn.Write([]byte("HTTP/1.1 404 OK\r\n\r\n"))
	}

	conn.Close()

}

func ParseHttpRequest(data []byte) Request {
	reqs := strings.Split(string(data), "\r\n")

	// make alocates the variable to memory
	var headers = make(Header)

	for _, v := range reqs {
		if strings.Contains(v, ":") {
			h := strings.Split(v, ":")
			if len(h) > 2 {
				// ??
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
	}
}
