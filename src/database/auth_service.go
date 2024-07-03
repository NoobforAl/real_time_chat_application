package database

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/NoobforAl/real_time_chat_application/src/entity"
	"github.com/NoobforAl/real_time_chat_application/src/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authService struct {
	client auth.AuthServiceClient
}

func newAuthService(grpcAddr string, log contract.Logger) *authService {
	// TODO: remove debug tls connection
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &authService{client: auth.NewAuthServiceClient(conn)}
}

func (s *store) Login(ctx context.Context, token string) (entity.User, error) {
	user, err := s.authClient.client.Login(ctx, &auth.LoginRequest{Token: token})
	if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		Id:           user.Id,
		Username:     user.Username,
		Notification: user.Notification,
	}, nil
}
