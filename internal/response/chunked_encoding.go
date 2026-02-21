package response

import (
	"fmt"
	"strconv"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerStatus != WriterBody {
		return 0, fmt.Errorf("cannot write body in state: %d", w.writerStatus)
	}
	datLen := len(p)
	hexLen := strconv.FormatInt(int64(datLen), 16)

	// written, err := w.writer.Write([]byte(fmt.Sprintf("%s\r\n%v\r\n", hexLen, p)))
	// writing data using separate lines so we don't have to string convert byte data
	totBytes := 0
	b, err := w.writer.Write([]byte(hexLen + "\r\n"))
	// can also do
	// b, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return 0, err
	}
	totBytes += b
	b, err = w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	totBytes += b
	b, _ = w.writer.Write([]byte("\r\n"))
	totBytes += b
	return totBytes, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	written, err := w.writer.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, err
	}
	return written, nil
}
