package server

import (
	"net/http"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

// NewServer creates a new Server for registering endpoints and
// starting the server.
func NewServer(addr string) (*Server, error) {
	mux := http.NewServeMux()
	serv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	s := &Server{
		httpServer: serv,
		mux:        mux,
	}

	return s, nil
}

// Start starts the server and listens on the address. It will return an
// error if any errors occur during the start up.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// RegisterHandler registers a new handler for the Server.
// This will panic if an existing pattern is given to register.
func (s *Server) RegisterHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}
