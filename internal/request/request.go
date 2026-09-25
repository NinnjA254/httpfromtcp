package request

// package request

import (
	// "errors"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/NinnjA254/httpfromtcp/internal/headers"
)

type parseState string

const (
	INIT    parseState = "init"
	HEADERS parseState = "headers"
	BODY    parseState = "body"
	DONE    parseState = "done"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parseState
}

func (r *Request) parse(data []byte) (int, error) {
	bytesParsed := 0

	for {
		currentlyParsing := data[bytesParsed:]
		// fmt.Printf("Parsing %q\n", currentlyParsing)
		switch r.state {
		case INIT:
			requestLine, n, err := parseRequestLine(currentlyParsing)

			if err != nil {
				return 0, err
			}
			if n == 0 {
				return bytesParsed, nil
			}
			r.RequestLine = *requestLine
			bytesParsed += n
			r.state = HEADERS
		case HEADERS:
			n, done, err := r.Headers.Parse(currentlyParsing)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				return bytesParsed, nil
			}
			bytesParsed += n
			if done {
				r.state = BODY
			}
		case BODY:
			contentLength := 0
			contentLengthString, ok := r.Headers.Get("content-length")
			if ok {
				n, err := strconv.Atoi(contentLengthString)
				if err != nil || n < 0 {
					return 0, errors.New("malformed content length")
				}
				contentLength = n
			}
			fmt.Printf("parsing body-> %q-> len %d\n", currentlyParsing, len(currentlyParsing))
			fmt.Println(" | - Content-Length:", contentLength)
			if len(currentlyParsing) < contentLength {
				return bytesParsed, nil
			}
			r.Body = make([]byte, contentLength)
			copy(r.Body, currentlyParsing[:contentLength])
			r.state = DONE
			return bytesParsed, nil
		}
	}
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

var ErrMalformedRequestLine = fmt.Errorf("malformed request line")

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLine := data[:idx]
	fmt.Printf("parsing requestLine-> %q-> len %d\n", requestLine, idx)
	parts := bytes.Split(requestLine, []byte(" "))
	if len(parts) != 3 {
		return nil, idx + 2, ErrMalformedRequestLine
	}
	method := string(parts[0])
	//method should have only capital alphabetic characters
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, idx + 2, ErrMalformedRequestLine
		}
	}
	requestTarget := string(parts[1])
	//requestTarget
	httpVersion := string(parts[2])
	if len(httpVersion) != 8 ||
		httpVersion[:5] != "HTTP/" || //test
		httpVersion[5] < '0' || httpVersion[5] > '9' ||
		httpVersion[6] != '.' ||
		httpVersion[7] < '0' || httpVersion[7] > '9' {
		return nil, idx + 2, ErrMalformedRequestLine
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion[5:],
	}, idx + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	dataLen := 0
	dataStart := 0
	buf := make([]byte, 1024)
	r := &Request{state: INIT, Headers: headers.NewHeaders()}

	for r.state != DONE {
		if dataLen == cap(buf) {
			buf = append(buf, 0)
			buf = buf[:cap(buf)]
			// fmt.Println("grow buf, cap now: ", cap(buf))
			// fmt.Println("===================")
		}
		fmt.Println()
		fmt.Println("Getting data from the connection...")
		n, err := reader.Read(buf[dataLen:])
		if err != nil {
			fmt.Println()
			return nil, err
		}
		dataLen += n

		data := buf[dataStart:dataLen]
		// fmt.Printf("data:-> %q\n", data)
		fmt.Printf("==================data================\n")
		fmt.Printf("%q\n", data)
		// fmt.Printf("%s\n", data)
		fmt.Printf("======================================\n")
		bytesParsed, err := r.parse(data)
		if err != nil {
			return nil, err
		}
		dataStart += bytesParsed
		// PrintRequest(r)
	}

	// PrintRequest(r)
	// fmt.Println()
	return r, nil
}
