package http

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
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

type Response struct {
	Protocol   string
	StatusCode int
	StatusText string
	Header     Header
	Body       []byte
}

func (R *Response) Write(conn net.Conn) {

	if R.Header != nil && R.Body == nil {

		conn.Write([]byte(fmt.Sprintf("%v %v %v\r\n%v\r\n\r\n", R.Protocol, R.StatusCode, R.StatusText, R.Header)))
	} else if R.Body != nil {

		conn.Write([]byte(fmt.Sprintf("%v %v %v\r\n%v\r\n\r\n%v", R.Protocol, R.StatusCode, R.StatusText, R.Header, string(R.Body))))
	} else {

		conn.Write([]byte(fmt.Sprintf("%v %v %v\r\n", R.Protocol, R.StatusCode, R.StatusText)))
	}

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

	req := parseHttpRequest(data)

	if req.Url == "/" {
		res := Response{
			Protocol:   "HTTP/1.1",
			StatusCode: 200,
			StatusText: "OK",
		}

		res.Write(conn)

	} else if strings.HasPrefix(req.Url, "/echo/") {

		vars, _ := strings.CutPrefix(req.Url, "/echo/")
		var headers = make(Header)
		headers["Content-Type"] = "text/plain"
		headers["Content-length"] = strconv.Itoa(len(vars))

		res := Response{
			Protocol:   "HTTP/1.1",
			StatusCode: 200,
			StatusText: "OK",
			Header:     headers,
			Body:       []byte(vars),
		}

		// conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(vars), vars)))
		res.Write(conn)

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

			_, err = f.Write(req.Body)
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

func parseHttpRequest(data []byte) Request {
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
		Body:   []byte(body),
	}
}
