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
	key, value := strings.TrimSuffix(fields[0], ":"), fields[1]
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

	// check to see if key already has values set
	tmp := h.Get(key)
	fmt.Printf("existing value: %s", tmp)
	if tmp == "" {
		h.Set(key, value)
	} else {
		concat := fmt.Sprintf("%s, %s", tmp, value)
		h.Set(key, concat)
	}

	return len(data[:idx+len(crlf)]), false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Get(key string) string {
	val, ok := h[strings.ToLower(key)]
	if ok {
		return val
	}
	return ""
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
