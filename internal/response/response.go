package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jwoodsiii/httpfromtcp/internal/headers"
)

type writerState int

const (
	WriterStatusLine writerState = iota
	WriterHeaders
	WriterBody
)

type Writer struct {
	writer       io.Writer
	writerStatus writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerStatus: WriterStatusLine,
		writer:       w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerStatus != WriterStatusLine {
		return fmt.Errorf("error: cannot write status line in state: %d", w.writerStatus)
	}
	defer func() { w.writerStatus = WriterHeaders }()
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

// TODO: Implement writeheaders on writer
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerStatus != WriterHeaders {
		return fmt.Errorf("error: cannot write headers in state: %d", w.writerStatus)
	}
	defer func() { w.writerStatus = WriterBody }()
	for k, v := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return fmt.Errorf("error writing header: %v", err)
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerStatus != WriterBody {
		return 0, fmt.Errorf("error: writer must be in body state")
	}
	return w.writer.Write(p)
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
