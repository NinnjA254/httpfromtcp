package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

// func TestRequestLineParse(t *testing.T) {}
func TestRequestFromReader(t *testing.T) {
	numBytesPerRead := 8

	// valid Request with no headers
	reader := chunkReader{
		data:            "GET / HTTP/1.1\r\n\r\n",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err := RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
	require.Equal(t, 0, len(r.Headers))

	// valid Request with one header
	reader = chunkReader{
		data: "GET / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n\r\n",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
	require.Equal(t, 1, len(r.Headers))
	require.Equal(t, "localhost:42069", r.Headers["host"])

	// valid Request with multiple headers
	reader = chunkReader{
		data: "GET / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Foo: bar\r\n\r\n",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
	require.Equal(t, 2, len(r.Headers))
	require.Equal(t, "localhost:42069", r.Headers["host"])
	require.Equal(t, "bar", r.Headers["foo"])

	// valid Request with multiple headers and a body
	//
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-type: application/json\r\n" +
			"content-length: 14\r\n" +
			"Accept: */*\r\n\r\n" +
			"{\"Foo\": \"bar\"}",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "POST", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
	require.Equal(t, 5, len(r.Headers))
	require.Equal(t, "*/*", r.Headers["accept"])
	require.Equal(t, "application/json", r.Headers["content-type"])
	require.Equal(t, "14", r.Headers["content-length"])
	require.Equal(t, "localhost:42069", r.Headers["host"])
	require.Equal(t, "curl/7.81.0", r.Headers["user-agent"])

	// valid Request with header-like data after blank line
	reader = chunkReader{
		data: "GET / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n\r\n" +
			"Not-Header: notheader\r\n\r\n",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
	require.Equal(t, 1, len(r.Headers))
	require.Equal(t, "localhost:42069", r.Headers["host"])
	_, ok := r.Headers["not-header"]
	require.False(t, ok)

	// Test: incomplete request
	reader = chunkReader{
		data: "GET / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"Accept:",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.Error(t, err)
	require.Nil(t, r)
}
func TestBodyParsing(t *testing.T) {
	numBytesPerRead := 3
	// content-length < actual content's length
	reader := chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-length: 2\r\n\r\n" +
			"123456789",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err := RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "12", string(r.Body))

	// content-length == actual content's length
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-length: 9\r\n\r\n" +
			"123456789",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "123456789", string(r.Body))

	// content-length > actual content's length
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-length: 10\r\n\r\n" +
			"123456789",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.Error(t, err)
	require.Equal(t, io.EOF, err)
	require.Nil(t, r)

	// content-length == 0
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-length: 0\r\n\r\n" +
			"123456789",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "", string(r.Body))

	// no content-length header
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n\r\n" +
			"123456789",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "", string(r.Body))

	// empty body, no content-length header
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n\r\n",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "", string(r.Body))

	// empty body, content-length == 0
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-length: 0\r\n\r\n",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "", string(r.Body))

	// content-length < 0
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-length: -1\r\n\r\n" +
			"123456789",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.Error(t, err)
	require.Nil(t, r)

	// decimal content-length
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-length: 1.2\r\n\r\n" +
			"123456789",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.Error(t, err)
	require.Nil(t, r)

	// non-digit content-length
	reader = chunkReader{
		data: "POST / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"content-length: foo1\r\n\r\n" +
			"123456789",
		numBytesPerRead: numBytesPerRead,
		pos:             0,
	}
	r, err = RequestFromReader(&reader)
	require.Error(t, err)
	require.Nil(t, r)
}
