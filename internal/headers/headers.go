package headers

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

var isFieldName = regexp.MustCompile("^[A-Za-z0-9!#$%&*+-.^_`|~]+$").MatchString

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// add newly parsed header kv pairs
	// return num bytes consumed, whether or not parsing is finished, err

	// look for crlf
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		log.Printf("no data consumed")
		return 0, false, nil
	} else if idx <= len(data[:2]) {
		log.Printf("found crlf at beginning of data")
		return 2, true, nil
	}
	strData := strings.TrimSpace(string(data))
	fields := strings.Fields(strData)
	key, value := strings.ToLower(strings.TrimSuffix(fields[0], ":")), fields[1]
	if strings.IndexByte(strData, ' ') <= len(key) {
		return 0, false, fmt.Errorf("invalid key format, no space allowed between host and colon")
	}
	//regexp impl of fieldname invalid char check
	// if !isFieldName(key) {
	// 	return 0, false, fmt.Errorf("invalid key format, key must match field-name spec")
	// }

	// loop impl of fieldname invalid char check
	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}

	//fmt.Printf("key- %s\nvalue- %s\n", key, value)
	h[key] = value
	return len(data[:idx+len(crlf)]), false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

// validTokens checks if the data contains only valid tokens
// or characters that are allowed in a token
func validTokens(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}
