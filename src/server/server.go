package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

type Server struct {
	Listener net.Listener
	Mux      *http.ServeMux
}

func NewServer(addr string) (*Server, error) {
	mux := http.NewServeMux()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Listener: l,
		Mux:      mux,
	}

	return s, nil
}

func (s *Server) CreateServer() *http.Server {
	s.Mux.HandleFunc("/", defaultRoot)

	ser := &http.Server{
		Addr:    s.Listener.Addr().String(),
		Handler: s.Mux,
	}

	return ser
}

func defaultRoot(w http.ResponseWriter, _ *http.Request) {
	fmt.Println(io.WriteString(w, "Hello!\n"))
	fmt.Println("Accessed")
}
