package httpserver

import (
	"net"
	"time"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = t
	}
}

func WriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = t
	}
}

func ShutdownTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = t
	}
}
