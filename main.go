package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		defer f.Close()
		var line string
		buf := make([]byte, 8, 8)
		for {
			b, err := f.Read(buf)
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				fmt.Printf("error reading bytes from file: %v", err)
				return
			}
			parts := strings.Split(string(buf[:b]), "\n")
			if len(parts) == 1 {
				line += parts[0]
				continue
			}
			for i := range parts[:len(parts)-1] {
				line += parts[i]
				out <- line
			}
			line = ""
			line += parts[len(parts)-1]
		}
		if line != "" {
			out <- line
		}
	}()
	return out
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	words := getLinesChannel(file)
	for l := range words {
		fmt.Printf("read: %s\n", l)
	}
}
