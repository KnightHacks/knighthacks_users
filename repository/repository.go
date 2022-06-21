package repository

import (
	"context"

	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_users/graph/model"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByOAuthUID(ctx context.Context, oAuthUID string, provider models.Provider) (*model.User, error)

	UpdateUser(ctx context.Context, id string, input model.NewUser) (*model.User, error)

	GetOAuth(ctx context.Context, userId string) (*model.OAuth, error)

	GetUsers(ctx context.Context, first int, after string) ([]*model.User, int, error)
	SearchUser(ctx context.Context, name string) ([]*model.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	CreateUser(ctx context.Context, oAuth *model.OAuth, input *model.NewUser) (*model.User, error)
}
