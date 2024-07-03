package contract

import (
	"context"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

type StoreUser interface {
	User(ctx context.Context, username string) (entity.User, error)
	UserIds(ctx context.Context) ([]string, error)
	CreateUser(ctx context.Context, userData entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, id string, updatedData entity.User) (entity.User, error)
}

type AuthenticationService interface {
	Login(ctx context.Context, token string) (entity.User, error)
}
