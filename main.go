package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer file.Close()
	buf := make([]byte, 8)
	for {
		b, err := file.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading bytes from file: %v", err)
		}
		fmt.Printf("read: %s\n", string(buf[:b]))
	}
}
