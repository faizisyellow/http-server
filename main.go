package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Header map[string]string

type Request struct {
	Method string
	Url    string
	Proto  string
	Header Header
	Body   string
}

func main() {
	ListenAndServe("4221")
}

func ListenAndServe(p string) {
	port := fmt.Sprintf("0.0.0.0:%v", p)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	handleConnection(conn)
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
	} else if strings.HasPrefix(req.Url, "/files/") {

		if req.Method == "POST" {
			filename, _ := strings.CutPrefix(req.Url, "/files/")
			path := fmt.Sprintf("files/%v", filename)
			f, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
				return
			}

			defer f.Close()

			_, err = f.Write([]byte(req.Body))
			if err != nil {
				log.Fatal(err)
			}

			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 201 Created\r\n\r\n")))
			return
		}

		dirpath, _ := strings.CutPrefix(req.Url, "/")
		var fileName string

		// if the directory empty it return dot
		if filepath.Dir(dirpath) != "." {
			paths, err := filepath.Glob(dirpath + ".*")
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(paths) == 0 {
				conn.Write([]byte("HTTP/1.1 404 OK\r\n\r\n"))
				return
			}

			fileName = paths[0]
		} else {
			conn.Write([]byte("HTTP/1.1 404 OK\r\n\r\n"))
			return
		}

		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n%v", len(string(fileContent)), string(fileContent))))

	} else {

		conn.Write([]byte("HTTP/1.1 404 OK\r\n\r\n"))
	}

	conn.Close()

}

func ParseHttpRequest(data []byte) Request {
	reqs := strings.Split(string(data), "\r\n")
	body := strings.Split(string(data), "\r\n\r\n")[1]

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
		Body:   body,
	}
}
