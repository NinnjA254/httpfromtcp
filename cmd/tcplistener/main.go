package main

import (
	"fmt"
	"log"
	"net"

	"github.com/NinnjA254/httpfromtcp/internal/request"
)

func main() {
	go ServeEcho(42068)
	// return

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error:", err)
		}
		// linesChannel := getLinesChannel(conn)

		fmt.Println()
		fmt.Printf("New request from->%v\n", conn.RemoteAddr().String())
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println()
			fmt.Println()
			fmt.Println("Error parsing request->", err)

			conn.Write([]byte("HTTP/1.1 400 Bad request\r\n"))
			conn.Close()
		} else {
			fmt.Println()
			fmt.Println()
			fmt.Println("Final request:")
			fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
				r.RequestLine.Method,
				r.RequestLine.RequestTarget,
				r.RequestLine.HttpVersion,
			)
			// fmt.Printf("Headers:\n- %v\n", r.Headers)
			fmt.Printf("Headers: (%d)\n", len(r.Headers))
			for k, v := range r.Headers {
				fmt.Printf("- %s:%s\n", k, v)
			}
			fmt.Println("Body:")
			fmt.Printf("%q\n", r.Body)

			conn.Write([]byte("HTTP/1.1 200 OK\r\n" +
				"Content-Type : application/json\r\n" +
				"Content-Length: 42\r\n\r\n" +
				"{\"message\":\"hello, world!\",\"status\":\"ok\"}\n"))
			conn.Close()
		}
	}
}
