package repository

import (
	"context"
	"github.com/KnightHacks/knighthacks_users/graph/model"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByAuthToken(ctx context.Context, authToken string) (*model.User, error)
	GetOAuth(ctx context.Context, userId string) (*model.OAuth, error)
}
