package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/NinnjA254/httpfromtcp/internal/request"
	"github.com/NinnjA254/httpfromtcp/internal/server"
)

var directory = flag.String("directory", ".", "the directory to serve")

func handler(req *request.Request, res *server.Response) {
	uri := req.RequestLine.RequestTarget
	filePath := *directory + uri
	fileInfo, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) == true || errors.Is(err, syscall.ENOTDIR) {
		res.Status = 404
		res.ReasonPhrase = "File not found"
		res.Headers.Set("content-type", "text/html; charset=utf-8")
		res.Body = []byte("<h1> 404 not found</h1>")
		res.Send()

		fmt.Println("not found:", filePath)
		return
	}
	if err != nil {
		res.Status = 500
		res.ReasonPhrase = "internal server error"
		res.Send()

		fmt.Println("internal server error:", err)
		return
	}
	if !fileInfo.IsDir() {
		f, err := os.Open(filePath)
		if err != nil {
			res.Status = 500
			res.ReasonPhrase = "internal server error"
			res.Send()

			fmt.Println("internal server error:", err)
			return
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			res.Status = 500
			res.ReasonPhrase = "internal server error"
			res.Send()

			fmt.Println("internal server error:", err)
			return
		}
		ext := filepath.Ext(filePath)
		contentType := mime.TypeByExtension(ext)
		contentLength := fmt.Sprintf("%d", len(data))
		res.Status = 200
		res.ReasonPhrase = "OK"
		res.Headers.Set("content-type", contentType)
		res.Headers.Set("content-length", contentLength)
		res.Body = data
		res.Send()
		return
	}

	if !strings.HasSuffix(filePath, "/") {
		res.Status = 301
		res.Headers.Set("location", uri+"/")
		res.Send()
		fmt.Printf("redirect %s -> %s\n", uri, uri+"/")

		return
	}
	entries, err := os.ReadDir(filePath)
	if err != nil {
		res.Status = 500
		res.ReasonPhrase = "internal server error"
		res.Send()

		fmt.Println("internal server error:", err)
		return
	}
	for _, entry := range entries {
		if entry.Name() == "index.html" && !entry.IsDir() {
			filePath += "index.html"
			f, err := os.Open(filePath)
			if err != nil {
				res.Status = 500
				res.ReasonPhrase = "internal server error"
				res.Send()

				fmt.Println("internal server error:", err)
				return
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				res.Status = 500
				res.ReasonPhrase = "internal server error"
				res.Send()

				fmt.Println("internal server error:", err)
				return
			}
			ext := filepath.Ext(filePath)
			contentType := mime.TypeByExtension(ext)
			contentLength := fmt.Sprintf("%d", len(data))
			res.Status = 200
			res.ReasonPhrase = "Ok"
			res.Headers.Set("content-type", contentType)
			res.Headers.Set("content-length", contentLength)
			res.Body = data
			res.Send()
			return
		}
	}

	//listing
	listing := "<body>"
	listing += fmt.Sprintf("<h1> Directory listing for %s</h1><hr/>", uri)
	listing += "<ul>"
	for _, entry := range entries {
		name := entry.Name()
		listing += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, uri+name, name)
	}
	listing += "</ul>"
	listing += "<hr/>"
	listing += "</body>"
	res.Status = 200
	res.ReasonPhrase = "OK"
	res.Headers.Set("content-type", "text/html")
	contentLength := fmt.Sprintf("%d", len(listing))
	res.Headers.Set("content-length", contentLength)
	res.Body = []byte(listing)
	res.Send()
}
func main() {
	flag.Parse()
	s := server.InitServer(handler)
	err := s.Listen(42068)
	if err != nil {
		fmt.Print("unexpected error starting server:", err)
	}
}
