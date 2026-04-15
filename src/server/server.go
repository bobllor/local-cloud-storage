package server

import (
	"net/http"
)

type Server struct {
	httpServer *http.Server
	Handler    *http.ServeMux
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
		Handler:    mux,
	}

	return s, nil
}

// Start starts the server and listens on the address. It will return an
// error if any errors occur during the start up.
//
// This should only used for development. Use s.StartTLS() for production.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// StartTLS starts the server and listens on the address with TLS. It will return
// an error if any errors occur during the start up.
//
// This requires the cert and key files.
func (s *Server) StartTLS(certFile string, keyFile string) error {
	// TODO: need to configure TLS here.
	// since this is a local project i think i can go with insecure and
	// bare minimum TLS settings, but ill look into the options.
	// most likely ill make this configurable via yaml.

	return s.httpServer.ListenAndServeTLS(certFile, keyFile)
}

// RegisterHandlerFunc registers a new handler function. It is not wrapped in middleware.
// This will panic if an existing pattern is given to register.
func (s *Server) RegisterHandlerFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.Handler.HandleFunc(pattern, handler)
}

// RegisterAuthHandler registers a new handler.
// This will panic if an existing pattern is given to register.
func (s *Server) RegisterHandler(pattern string, handler http.Handler) {
	s.Handler.Handle(pattern, handler)
}
