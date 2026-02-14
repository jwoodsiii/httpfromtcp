package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(raw string) (*RequestLine, error) {
	// fmt.Printf("Raw string: %q\n", raw)
	rawRl := strings.Fields(strings.Split(raw, "\r\n")[0])
	if len(rawRl) != 3 {
		return nil, fmt.Errorf("malformed request")
	}

	// validate parsed information before building struct
	// Method only has capital alphabetic characters
	method := rawRl[0]
	// no whitespace allowed in request target
	target := strings.TrimSpace(rawRl[1])
	version := strings.Split(rawRl[2], "/")

	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	if version[1] != "1.1" {
		return nil, fmt.Errorf("incorrect http version: %s", version[1])
	} else if version[0] != "HTTP" {
		return nil, fmt.Errorf("invalid http version")
	}
	rl := &RequestLine{
		HttpVersion:   version[1],
		RequestTarget: target,
		Method:        method,
	}
	return rl, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading string request: %v", err)
	}
	rl, err := parseRequestLine(string(req))
	if err != nil {
		return nil, fmt.Errorf("error parsing string to requestline: %v", err)
	}

	return &Request{RequestLine: *rl}, nil
}
