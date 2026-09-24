package main

import (
	"fmt"
	"io"
	"net/http"
)

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Println()
	fmt.Printf("New %v request from: %v\n", r.Method, r.RemoteAddr)
	fmt.Printf("Headers (%v)\n", len(r.Header))
	for k, v := range r.Header {
		fmt.Printf("- %q: %q\n", k, v)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error getting body:", err)
		return
	}
	fmt.Println("Body")
	fmt.Printf("- %q\n", body)
	w.Write(body)
	// io.Copy(w, body)
}

func ServeEcho(port int) {
	http.HandleFunc("/", echo)

	fmt.Println("echo server up!")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
