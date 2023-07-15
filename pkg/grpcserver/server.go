package grpcserver

import (
	"google.golang.org/grpc"
	"log"
	"net"
)

const _defaultAddr = ":50051"

var RpcServer *grpc.Server

type Server struct {
	server *grpc.Server
	notify chan error
	Addr   string
}

func New(opts ...Option) *Server {
	grpcServer := grpc.NewServer()
	RpcServer = grpcServer
	s := &Server{
		server: grpcServer,
		notify: make(chan error, 1),
		Addr:   _defaultAddr,
	}

	for _, opt := range opts {
		opt(s)
	}
	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		lis, err := net.Listen("tcp", s.Addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s.notify <- s.server.Serve(lis)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() {
	s.server.GracefulStop()
}
