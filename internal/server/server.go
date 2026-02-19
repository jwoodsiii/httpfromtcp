package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	response "github.com/jwoodsiii/httpfromtcp/internal/repsonse"
	"github.com/jwoodsiii/httpfromtcp/internal/request"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	state    *atomic.Bool // true if server is running
	Listener net.Listener
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Error      string
}

func Serve(port int, handler Handler) (*Server, error) {
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
		handler:  handler,
	}
	go s.listen()
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
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusServerError,
			Error:      err.Error(),
		}
		hErr.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}
	b := buf.Bytes()
	// update to use internal/response to format responses
	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		fmt.Printf("error writing status line: %v", err)
		return
	}
	headers := response.GetDefaultHeaders(len(b))

	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("error writing headers: %v", err)
		return
	}

	conn.Write(b)
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	msgBytes := []byte(he.Error)
	headers := response.GetDefaultHeaders(len(msgBytes))
	response.WriteHeaders(w, headers)
	w.Write(msgBytes)
}
