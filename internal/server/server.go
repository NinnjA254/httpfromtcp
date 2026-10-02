package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"

	"github.com/NinnjA254/httpfromtcp/cli"
	"github.com/NinnjA254/httpfromtcp/internal/request"
)

type routeHandler func(req *request.Request, res *Response)
type Server struct {
	routeHandlers   map[string]routeHandler
	notfoundHandler routeHandler
}

func InitServer() Server {
	flag.Parse()
	return Server{
		routeHandlers:   map[string]routeHandler{},
		notfoundHandler: defaultnotfoundHandler,
	}
}
func (s *Server) Route(path string, handler routeHandler) {
	//add a check to warn if a route is overwritten
	s.routeHandlers[path] = handler
}
func defaultnotfoundHandler(req *request.Request, res *Response) {
	fmt.Println("404 not found:", req.RequestLine.RequestTarget)
	res.SendText("not found, please stand by, proper 404 messages under construction")
}
func (s *Server) Listen(port int) (err error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Printf("Server up! Listening on port %d\n", port)
	fmt.Printf("Currently registered routes:\n")
	for k := range s.routeHandlers {
		fmt.Printf("    %v\n", k)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			cli.Printf("Error accepting request, %v", err)
			continue
		}
		handleRequest(s, conn)
	}
}
func handleRequest(s *Server, conn net.Conn) {
	defer conn.Close()
	fmt.Printf("\nNew request from->%v\n", conn.RemoteAddr().String())
	req, err := request.RequestFromReader(conn)
	if err != nil {
		//send error http response
		return
	}
	requestTarget := req.RequestLine.RequestTarget
	handler, ok := s.routeHandlers[requestTarget]
	if !ok {
		handler = s.notfoundHandler
	}
	res := &Response{conn: conn}
	handler(req, res)
}

type Response struct {
	conn net.Conn
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
	contentType := "text/html"
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
