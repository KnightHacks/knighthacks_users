package repository

import (
	"context"
	"github.com/KnightHacks/knighthacks_users/graph/model"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByOAuthUID(ctx context.Context, oAuthUID string) (*model.User, error)
	GetOAuth(ctx context.Context, userId string) (*model.OAuth, error)

	CreateUser(ctx context.Context, oAuth *model.OAuth, input *model.NewUser) (*model.User, error)
}
