//go:generate protoc -I=../proto/ --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../proto/auth.proto
package auth

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/services/auth/jwt"
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

// TODO: need database ??
func (s *server) Login(ctx context.Context, in *LoginRequest) (*LoginInfoReply, error) {
	info_token, err := jwt.ValidateToken(in.Token, []byte(config.SecretKey()))
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return &LoginInfoReply{
		Id:           info_token.Id,
		Username:     info_token.Username,
		Notification: info_token.Notifications,
	}, nil
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
