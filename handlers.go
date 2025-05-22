package main

import (
	"github.com/faizisyellow/http-server/http"
)

func base(r http.Request, w *http.ResponseWrite) {
	w.Write(200, "OK")
}

// func echo(r http.Request, w *http.ResponseWritter) {

// }

// func userAgent(r http.Request, w *http.ResponseWritter) {

// }

// func fileOperation(r http.Request, w *http.Response) {
// 	if r.Method == "POST" {
// 		filename, _ := strings.CutPrefix(r.Url, "/files/")
// 		path := fmt.Sprintf("files/%v", filename)
// 		f, err := os.Create(path)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}

// 		defer f.Close()

// 		_, err = f.Write(r.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 201 Created\r\n\r\n")))
// 		return
// 	}

// 	dirpath, _ := strings.CutPrefix(r.Url, "/")
// 	var fileName string

// 	// if the directory empty it return dot
// 	if filepath.Dir(dirpath) != "." {
// 		paths, err := filepath.Glob(dirpath + ".*")
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		if len(paths) == 0 {
// 			conn.Write([]byte("HTTP/1.1 404 OK\r\n\r\n"))
// 			return
// 		}

// 		fileName = paths[0]
// 	} else {
// 		conn.Write([]byte("HTTP/1.1 404 OK\r\n\r\n"))
// 		return
// 	}

// 	fileContent, err := os.ReadFile(fileName)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n%v", len(string(fileContent)), string(fileContent))))

// }
