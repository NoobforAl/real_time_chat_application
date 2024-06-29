//go:generate protoc -I=../proto/ --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../proto/auth.proto
package auth

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	grpc "google.golang.org/grpc"
)

type server struct {
	UnsafeAuthServiceServer
	store contract.Store
	log   contract.Logger
}

var onc sync.Once
var grpc_server *server

func New(store contract.Store, logger contract.Logger) *server {
	onc.Do(func() {
		grpc_server = &server{
			store: store,
			log:   logger,
		}
	})

	return grpc_server
}

func (s *server) Login(ctx context.Context, in *LoginRequest) (*LoginInfoReply, error) {
	return &LoginInfoReply{}, nil
}

func (s *server) Run(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	RegisterAuthServiceServer(server, s)

	s.log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		s.log.Fatalf("failed to serve: %v", err)
	}
}
