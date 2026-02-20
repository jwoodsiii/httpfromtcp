package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	StatusOK          StatusCode = 200
	StatusClientError StatusCode = 400
	StatusServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusClientError:
		reasonPhrase = "Bad Request"
	case StatusServerError:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))
	return err
}
