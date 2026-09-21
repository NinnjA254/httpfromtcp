package main

import (
	"fmt"
	"io"
	"net/http"
)

func echo(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%q\n", r.Header)
	for k, v := range r.Header {
		fmt.Printf("%q, %q\n", k, v)
	}
	io.Copy(w, r.Body)
}

func ServeEcho(port int) {
	http.HandleFunc("/", echo)

	fmt.Println("echo server up!")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
