package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/NinnjA254/httpfromtcp/cli"
	"github.com/NinnjA254/httpfromtcp/internal/headers"
	"github.com/NinnjA254/httpfromtcp/internal/request"
)

type requestHandler func(req *request.Request, res *Response)
type Server struct {
	handler requestHandler
}

func InitServer(h requestHandler) Server {
	flag.Parse()
	return Server{
		handler: h,
	}
}

func (s *Server) Listen(port int) (err error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Printf("Server up! Listening on port %d\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting request, %v", err)
			continue
		}
		handleRequest(s, conn)
	}
}
func handleRequest(s *Server, conn net.Conn) {
	defer conn.Close()
	cli.Printf("\nNew request from->%v\n", conn.RemoteAddr().String())
	req, err := request.RequestFromReader(conn)
	if err != nil {
		//send error http response
		return
	}
	if req.RequestLine.HttpVersion != "1.1" {
		//send error http response
		return
	}
	res := &Response{
		conn:        conn,
		httpVersion: req.RequestLine.HttpVersion,
		Headers: headers.Headers{
			"date":   time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"),
			"server": "mninjnx",
		},
		Body: nil,
	}
	s.handler(req, res)
}

type Response struct {
	conn         net.Conn
	httpVersion  string
	Status       int
	ReasonPhrase string
	Headers      headers.Headers
	Body         []byte
}

func (res *Response) Send() {
	responseLine := fmt.Sprintf("HTTP/%s %d %s\r\n", res.httpVersion, res.Status, res.ReasonPhrase)
	headers := ""
	for k, v := range res.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	headers += "\r\n"
	httpResponse := append([]byte(responseLine+headers), res.Body...)
	res.conn.Write([]byte(httpResponse))
}
func (res *Response) SendText(content string) {
	contentLength := len(content)
	contentType := "text/plain"
	serverName := "mninjnx"
	httpResponse := []byte("HTTP/1.1 200 OK\r\n" +
		fmt.Sprintf("Content-Type: %s\r\n", contentType) +
		fmt.Sprintf("Content-Length: %d\r\n", contentLength) +
		fmt.Sprintf("Server: %s\r\n", serverName) +
		"\r\n")
	httpResponse = append(httpResponse, content...)
	res.conn.Write([]byte(httpResponse))

}
func (res *Response) SendHTML(content string) {
	contentLength := len(content)
	contentType := "text/html; charset=UTF-8"
	serverName := "Mninjnx"
	httpResponse := []byte("HTTP/1.1 200 OK\r\n" +
		fmt.Sprintf("Content-Type: %s\r\n", contentType) +
		fmt.Sprintf("Content-Length: %d\r\n", contentLength) +
		fmt.Sprintf("Server: %s\r\n", serverName) +
		"\r\n")
	httpResponse = append(httpResponse, content...)
	res.conn.Write([]byte(httpResponse))
}
func (res *Response) SendJSON(content any) {
	data, err := json.Marshal(content)
	if err != nil {
		data = []byte(`{"error": "error marshalling json, please stand by"}`)
	}
	contentLength := len(data)
	contentType := "application/json"
	serverName := "Mninjnx"
	httpResponse := []byte("HTTP/1.1 200 OK\r\n" +
		fmt.Sprintf("Content-Type: %s\r\n", contentType) +
		fmt.Sprintf("Content-Length: %d\r\n", contentLength) +
		fmt.Sprintf("Server: %s\r\n", serverName) +
		"\r\n")
	httpResponse = append(httpResponse, data...)
	res.conn.Write([]byte(httpResponse))
}
func (res *Response) SendImage(imageData []byte) {
	contentLength := len(imageData)
	contentType := "image/jpeg"
	serverName := "Mninjnx"
	httpResponse := []byte("HTTP/1.1 200 OK\r\n" +
		fmt.Sprintf("Content-Type: %s\r\n", contentType) +
		fmt.Sprintf("Content-Length: %d\r\n", contentLength) +
		fmt.Sprintf("Server: %s\r\n", serverName) +
		"\r\n")
	httpResponse = append(httpResponse, imageData...)
	res.conn.Write([]byte(httpResponse))
}
