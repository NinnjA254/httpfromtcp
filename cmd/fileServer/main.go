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
	"time"

	"github.com/NinnjA254/httpfromtcp/internal/request"
	"github.com/NinnjA254/httpfromtcp/internal/server"
)

var directory = flag.String("directory", ".", "the directory to serve")

func handler(req *request.Request, res *server.Response) {
	if req.RequestLine.Method != "GET" {
		res.Status = 501
		res.ReasonPhrase = fmt.Sprintf("Unsupported Method (%s)", req.RequestLine.Method)
		res.Headers.Add("content-length", "0")
		res.Send()

		fmt.Printf("[%s] %s %s - %d %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), req.RequestLine.Method, req.RequestLine.RequestTarget, res.Status, res.ReasonPhrase)
		return
	}
	uri := req.RequestLine.RequestTarget
	filePath := *directory + uri
	fileInfo, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) == true || errors.Is(err, syscall.ENOTDIR) {
		res.Status = 404
		res.ReasonPhrase = "File not found"
		res.Headers.Add("content-type", "text/html; charset=utf-8")
		content := []byte("<h1> 404 not found</h1>")
		contentLength := fmt.Sprintf("%d", len(content))
		res.Headers.Add("content-length", contentLength)
		res.Body = content
		res.Send()

		fmt.Printf("[%s] %s %s - %d %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), req.RequestLine.Method, req.RequestLine.RequestTarget, res.Status, res.ReasonPhrase)
		return
	}
	if err != nil {
		res.Status = 500
		res.ReasonPhrase = "internal server error"
		res.Headers.Add("content-length", "0")
		res.Send()

		fmt.Println("internal server error:", err)
		return
	}
	if !fileInfo.IsDir() {
		f, err := os.Open(filePath)
		if err != nil {
			res.Status = 500
			res.Headers.Add("content-length", "0")
			res.ReasonPhrase = "internal server error"
			res.Send()

			fmt.Println("internal server error:", err)
			return
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			res.Status = 500
			res.Headers.Add("content-length", "0")
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
		res.Headers.Add("content-type", contentType)
		res.Headers.Add("content-length", contentLength)
		res.Body = data
		res.Send()

		fmt.Printf("[%s] %s %s - %d %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), req.RequestLine.Method, req.RequestLine.RequestTarget, res.Status, res.ReasonPhrase)
		return
	}

	if !strings.HasSuffix(filePath, "/") {
		res.Status = 301
		res.ReasonPhrase = "Moved Permanently"
		res.Headers.Add("location", uri+"/")
		res.Headers.Add("content-length", "0")
		res.Send()

		fmt.Printf("[%s] %s %s - %d %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), req.RequestLine.Method, req.RequestLine.RequestTarget, res.Status, res.ReasonPhrase)
		return
	}
	entries, err := os.ReadDir(filePath)
	if err != nil {
		res.Status = 500
		res.ReasonPhrase = "internal server error"
		res.Headers.Add("content-length", "0")
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
				res.Headers.Add("content-length", "0")
				res.Send()

				fmt.Println("internal server error:", err)
				return
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				res.Status = 500
				res.ReasonPhrase = "internal server error"
				res.Headers.Add("content-length", "0")
				res.Send()

				fmt.Println("internal server error:", err)
				return
			}
			ext := filepath.Ext(filePath)
			contentType := mime.TypeByExtension(ext)
			contentLength := fmt.Sprintf("%d", len(data))
			res.Status = 200
			res.ReasonPhrase = "OK"
			res.Headers.Add("content-type", contentType)
			res.Headers.Add("content-length", contentLength)
			res.Body = data
			res.Send()

			fmt.Printf("[%s] %s %s - %d %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), req.RequestLine.Method, req.RequestLine.RequestTarget, res.Status, res.ReasonPhrase)
			return
		}
	}

	//listing
	listing := "<body>"
	listing += fmt.Sprintf("<h1> Directory listing for %s</h1><hr/>", uri)
	listing += "<ul>"
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		listing += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, uri+name, name)
	}
	listing += "</ul>"
	listing += "<hr/>"
	listing += "</body>"
	res.Status = 200
	res.ReasonPhrase = "OK"
	res.Headers.Add("content-type", "text/html")
	contentLength := fmt.Sprintf("%d", len(listing))
	res.Headers.Add("content-length", contentLength)
	res.Body = []byte(listing)
	res.Send()

	fmt.Printf("[%s] %s %s - %d %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), req.RequestLine.Method, req.RequestLine.RequestTarget, res.Status, res.ReasonPhrase)
}
func main() {
	flag.Parse()
	s := server.InitServer(handler)
	err := s.Listen(42069)
	if err != nil {
		fmt.Println("unexpected error starting server:", err)
	}
}
