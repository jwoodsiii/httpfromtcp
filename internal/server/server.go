package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/jwoodsiii/httpfromtcp/internal/request"
	"github.com/jwoodsiii/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

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
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusClientError)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	msgBytes := []byte(he.Error)
	headers := response.GetDefaultHeaders(len(msgBytes))
	response.WriteHeaders(w, headers)
	w.Write(msgBytes)
}
