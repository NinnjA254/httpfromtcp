package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	buf := make([]byte, 8)
	str := ""
	linesChannel := make(chan string)
	go func() {
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
					linesChannel <- str
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
			linesChannel <- str
		}
		close(linesChannel)
	}()
	return linesChannel
}
func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("error opening file", err)
	}
	linesChannel := getLinesChannel(f)
	for line := range linesChannel {
		fmt.Printf("read: %s\n", line)
	}
}
