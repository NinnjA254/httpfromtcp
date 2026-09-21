package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

var ErrMalformedHeader = fmt.Errorf("malformed header")
var ErrMalformedHeaderName = fmt.Errorf("malformed header name")
var ErrMalformedHeaderValue = fmt.Errorf("malformed header value") //

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}
	fmt.Printf("parsing header-> %q-> len %d\n", string(data[:idx]), idx)
	if idx == 0 {
		return 2, true, nil
	}
	header := data[:idx]
	parts := bytes.SplitN(header, []byte(":"), 2)
	// parts = bytes.Split(header, []byte(":"))
	if len(parts) != 2 {
		return 0, false, ErrMalformedHeader
	}
	key := string(parts[0])
	val := strings.Trim(string(parts[1]), " \t")
	if !isValidFieldName(key) {
		fmt.Println("invalid field name")
		return 0, false, ErrMalformedHeaderName
	}
	if !isValidFieldValue(val) {
		return 0, false, ErrMalformedHeaderValue
	}

	key = strings.ToLower(key)
	if _, exists := h[key]; exists {
		h[key] += "," + val
		//
	} else {
		h[key] = strings.Trim(val, " ")
	}

	return idx + 2, false, nil
}

func isValidFieldName(s string) bool {
	for _, c := range s {
		if c < '!' || c > '~' { //range of US-ASCII visual characters
			return false
		}
	}
	// field name can't contain delimiters
	if strings.ContainsAny(s, "\"(),/:;<=>?@[\\]{}") {
		return false
	}
	return true
}
func isValidFieldValue(s string) bool {
	length := len(s)
	for idx, c := range s {
		if idx == 0 || idx == length-1 {
			if c < '!' || c > '~' {
				return false
			}
		} else if (c < '!' || c > '~') &&
			c != ' ' &&
			c != '\t' { //range of US-ASCII visual characters
			return false
		}
	}
	return true
}
