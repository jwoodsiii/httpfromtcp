package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
	lst, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error creating tcp listener: %v", err)
	}

	conn, err := lst.Accept()
	if err != nil {
		fmt.Printf("error creating connection: %v", err)
	}
	defer conn.Close()
	fmt.Printf("conn accepted...\n")
	words := getLinesChannel(conn)
	for l := range words {
		fmt.Printf("%s\n", l)
	}
	fmt.Printf("closing conn")
}
