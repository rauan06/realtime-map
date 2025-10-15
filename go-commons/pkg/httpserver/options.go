package httpserver

import (
	"net"
	"net/http"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

func Mux(mux *http.ServeMux) Option {
	return func(s *Server) {
		s.mux = mux
	}
}
