package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	response "github.com/jwoodsiii/httpfromtcp/internal/repsonse"
)

type Server struct {
	state    *atomic.Bool // true if server is running
	Listener net.Listener
}

func Serve(port int) (*Server, error) {
	// create net.listener and return server instance
	// start listening for reqs inside goroutine
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error creating listener: %v", err)
	}
	var serverState atomic.Bool
	serverState.Store(true)
	s := &Server{
		Listener: conn,
		state:    &serverState,
	}
	go func() {
		s.listen()
	}()

	return s, nil
}

func (s *Server) Close() error {
	// close listener and server
	s.state.Store(false)
	if err := s.Listener.Close(); err != nil {
		return fmt.Errorf("error attempting to close listener: %v", err)
	}
	return nil
}

func (s *Server) listen() {
	for {
		if !s.state.Load() {
			log.Printf("server closed")
			return
		}
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	// handle single conn by writing output and then closing conn
	// _, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!"))
	// if err != nil {
	// 	fmt.Printf("error writing data to connection: %v", err)
	// }

	// update to use internal/response to format responses
	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		fmt.Printf("error writing status line: %v", err)
		return
	}
	headers := response.GetDefaultHeaders(0)

	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("error writing headers: %v", err)
		return
	}

	if err := conn.Close(); err != nil {
		return
	}
}
