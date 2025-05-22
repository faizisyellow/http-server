package http

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type ServerResponse interface {
	Write(statusCode int, statusText string, headers Header, body string)
}

type Header map[string]string

type Request struct {
	Method string
	Url    string
	Proto  string
	Header Header
	Body   []byte
	Params map[string]string
}

type ResponseWrite struct {
	Conn net.Conn
}

func (r *ResponseWrite) Write(statusCode int, statusText string, headers Header, body string) {
	response := fmt.Sprintf("HTTP/1.1 %v %v\r\n", statusCode, statusText)

	if headers == nil {
		headers = make(Header)
	}

	//  default headers
	headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "text/plain"
	}

	for k, v := range headers {
		response += fmt.Sprintf("%v: %v\r\n", k, v)
	}

	response += "\r\n"
	response += body

	r.Conn.Write([]byte(response))
}

type HandlerFunc func(Request, ServerResponse)

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
		log.Fatal(err)
		return
	}

	req := parseHttpRequest(data[:n])

	for pattern, handler := range routes {
		if ok, params := matchRoute(pattern, req.Url); ok {
			req.Params = params
			handler(req, &ResponseWrite{Conn: conn})
			break
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

func matchRoute(pattern, path string) (bool, map[string]string) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	params := map[string]string{}
	for i := range patternParts {
		if i >= len(pathParts) {
			return false, nil
		}

		pp := patternParts[i]
		cp := pathParts[i]

		if strings.HasPrefix(pp, "*") {
			key := pp[1:]
			params[key] = strings.Join(pathParts[i:], "/")
			return true, params
		}

		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			key := pp[1 : len(pp)-1]
			params[key] = cp
		} else if pp != cp {
			return false, nil
		}
	}

	if len(pathParts) > len(patternParts) {
		return false, nil
	}

	return true, params
}
