package response

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/jwoodsiii/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK          = 200
	StatusClientError = 400
	StatusServerError = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	// map status code to correct reason phrase, we only support 3
	switch statusCode {
	case StatusOK:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusClientError:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusServerError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		log.Printf("empty reason code")
		w.Write([]byte("\r\n"))
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		h := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(h))
		if err != nil {
			return fmt.Errorf("error writing header: %v", err)
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
