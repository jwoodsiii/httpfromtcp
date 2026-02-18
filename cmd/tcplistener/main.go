package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jwoodsiii/httpfromtcp/internal/request"
)

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
	_, err = request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("error getting request from connection: %v", err)
	}
	// fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
	fmt.Printf("closing conn")
}
