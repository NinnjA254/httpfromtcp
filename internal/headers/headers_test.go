package headers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {

	// Valid single header without OWS (optional white space)
	headers := NewHeaders()
	data := []byte("Host:localhost:42069\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "localhost:42069", headers["host"])
	require.Equal(t, 22, n)
	require.False(t, done)

	// Valid single header with OWS
	headers = NewHeaders()
	data = []byte("Host:       localhost:42069       \r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "localhost:42069", headers["host"])
	require.Equal(t, 36, n)
	require.False(t, done)

	// Valid single pascal case header
	headers = NewHeaders()
	data = []byte("Accept-Encoding: deflate\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "deflate", headers["accept-encoding"])
	require.Equal(t, 26, n)
	require.False(t, done)
	_, ok := headers["Accept-Encoding"]
	require.False(t, done)

	// Valid single header with SP and HTAB within the header value
	headers = NewHeaders()
	data = []byte("Foo:hello \tworld\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "hello \tworld", headers["foo"])
	require.Equal(t, 18, n)
	require.False(t, done)

	// Valid blank header
	headers = NewHeaders()
	data = []byte("Host:\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "", headers["host"])
	require.Equal(t, 7, n)
	require.False(t, done)

	// Blank line(a blank line is the end of headers)
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.True(t, done)

	// Blank line + more data
	// ie. end of headers plus some data from the request body
	headers = NewHeaders()
	data = []byte("\r\n{\"name\": \"Mutangubia\"}\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.True(t, done)

	//headers that already exist
	// Valid header that already exists
	headers = NewHeaders()
	headers["accept-encoding"] = "gzip"
	data = []byte("accept-encoding: deflate\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "gzip,deflate", headers["accept-encoding"])
	require.Equal(t, 26, n)
	require.False(t, done)

	// Valid pascal case header that already exists
	headers = NewHeaders()
	headers["accept-encoding"] = "gzip"
	data = []byte("Accept-Encoding: deflate\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "gzip,deflate", headers["accept-encoding"])
	require.Equal(t, 26, n)
	require.False(t, done)

	// Valid header that already exists with more than one value
	headers = NewHeaders()
	headers["accept-encoding"] = "gzip,br,zstd"
	data = []byte("accept-encoding: deflate\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "gzip,br,zstd,deflate", headers["accept-encoding"])
	require.Equal(t, 26, n)
	require.False(t, done)

	// Valid blank header that already exists
	headers = NewHeaders()
	headers["accept-encoding"] = "gzip,br,zstd"
	data = []byte("accept-encoding:\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "gzip,br,zstd,", headers["accept-encoding"])
	require.Equal(t, 18, n)
	require.False(t, done)

	// Invalid header , missing colon
	headers = NewHeaders()
	data = []byte("Hostlocalhost\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
	_, ok = headers["host"]
	require.False(t, ok)

	// Invalid header name, leading white space in header name
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
	_, ok = headers["host"]
	require.False(t, ok)

	// Invalid header name, white space between header name and colon
	headers = NewHeaders()
	data = []byte("Host : localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
	_, ok = headers["host"]
	require.False(t, ok)

	// Invalid header name, illegal characters in header name // make a function that generates different illegal headers and tests them in a loop?
	headers = NewHeaders()
	data = []byte("Ho(st: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
	_, ok = headers["ho(st"]
	require.False(t, ok)
	// Invalid header name, illegal characters in header name // make a function that generates different illegal headers and tests them in a loop?
	headers = NewHeaders()
	data = []byte("Ho©st: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
	_, ok = headers["ho(st"]
	require.False(t, ok)

	//invalid header value, illegal characters in header value
	headers = NewHeaders()
	data = []byte("Host: local\vhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
	_, ok = headers["host"]
	require.False(t, ok)

}
