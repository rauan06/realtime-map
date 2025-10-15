// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultAddr         = ":80"
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second
)

type Server struct {
	App *http.Server

	mux          *http.ServeMux
	notify       chan error
	address      string
	writeTimeout time.Duration
	readTimeout  time.Duration
}

func New(opts ...Option) *Server {
	server := &Server{
		App: nil,

		mux:     http.NewServeMux(),
		notify:  make(chan error, 1),
		address: _defaultAddr,
	}

	for _, opt := range opts {
		opt(server)
	}

	app := &http.Server{
		Addr:    server.address,
		Handler: server.mux,

		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}

	server.App = app

	return server
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.App.ListenAndServe()

		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.App.Shutdown(context.Background())
}
