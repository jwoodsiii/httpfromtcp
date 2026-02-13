package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer file.Close()
	buf := make([]byte, 8, 8)
	var line string
	for {
		b, err := file.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			fmt.Printf("error reading bytes from file: %v", err)
		}
		parts := strings.Split(string(buf[:b]), "\n")

		if len(parts) == 1 {
			line += parts[0]
			continue
		}
		for i := range parts[:len(parts)-1] {
			line += parts[i]
			fmt.Printf("read: %s\n", line)
		}
		line = ""
		line += parts[len(parts)-1]
	}
	if line != "" {
		fmt.Printf("read: %s\n", line)
	}

}
