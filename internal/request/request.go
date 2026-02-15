package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type parseState int

const (
	initialized parseState = iota
	done
)

const crlf = "\r\n"
const bufSize = 8

type Request struct {
	RequestLine RequestLine
	state       parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// don't read all bytes at once, use loop+parse to continually pull from reader and parse
	// data in chunks
	// req, err := io.ReadAll(reader)
	buf := make([]byte, bufSize)
	readToIdx := 0 // track how much data we've read from reader into buffer
	r := &Request{state: initialized}
	for r.state != done {
		if readToIdx >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		bytesRead, err := reader.Read(buf[readToIdx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.state = done
				break
			} else {
				return nil, fmt.Errorf("error reading string request: %v", err)
			}
		}
		readToIdx += bytesRead                  // mv readidx to reflect new data
		parsed, err := r.parse(buf[:readToIdx]) // parse data we've placed in buffer
		if err != nil {
			return nil, fmt.Errorf("error parsing data from bytes: %v", err)
		}
		// remove parsed data from buffer
		copy(buf, buf[parsed:])
		readToIdx -= parsed
	}
	// rl, _, err := parseRequestLine(string(buf))
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing string to requestline: %v", err)
	// }
	// r.RequestLine = *rl
	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	// accept all unparsed bytes from buffer
	// update state of parser, updates request line struct
	// return num bytes successfully parsed and error, if one is present
	switch r.state {
	case initialized:
		rl, parsed, err := parseRequestLine(data)
		if err != nil {
			return 0, fmt.Errorf("error parsing request line from buffer: %v", err)
		}
		if parsed == 0 {
			fmt.Printf("need more data\n")
			return 0, nil
		}
		r.RequestLine = *rl
		r.state = done
		return parsed, nil
	case done:
		return 0, fmt.Errorf("error: attempting to read data in done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// fmt.Printf("Raw string: %q\n", raw)
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
