package request

import "fmt"

func PrintRequest(r *Request) {
	fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
		r.RequestLine.Method,
		r.RequestLine.RequestTarget,
		r.RequestLine.HttpVersion,
	)
	fmt.Println("Headers:")
	for k, v := range r.Headers {
		fmt.Printf("- %s:%s\n", k, v)
	}
	fmt.Println("Body:")
	fmt.Printf("%q\n", r.Body)
	// fmt.Println("Body:")
	// fmt.Println("work on body printing")
}
