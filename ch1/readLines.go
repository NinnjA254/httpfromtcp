package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

/*
chapter 1, lesson 4
read and print lines from a file
*/
func _main() { /*added underscore before main to silence "main redeclared" temporarily*/
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("error opening file", err)
	}
	buf := make([]byte, 8)
	str := ""
	for {
		n, err := f.Read(buf)
		defer f.Close()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("an unexpected error occured: %s\n", err)
			}
			break
		}
		// fmt.Printf("buf: %q %d bytes\n", string(buf[:n]), n)
		lines := strings.Split(string(buf[:n]), "\n")
		if len(lines) > 1 {
			// fmt.Printf("lines in the buffer: %q\n", lines)
			for _, line := range lines[:len(lines)-1] {
				str += line
				fmt.Printf("read: %s\n", str)
				str = ""
			}
			// fmt.Printf("left in the buffer: %q\n", lines[len(lines)-1])
			str = lines[len(lines)-1]
		} else {
			// fmt.Println("no lines in the buffer.")
			// str += string(buf[:n])
			str += lines[0]
		}
		// fmt.Println()
	}
	if str != "" {
		fmt.Printf("read: %s\n", str)
	}
}
