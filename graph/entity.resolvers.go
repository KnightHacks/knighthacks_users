package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/KnightHacks/knighthacks_users/graph/generated"
	"github.com/KnightHacks/knighthacks_users/graph/model"
)

func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := r.Resolver.Repository.GetUserByID(ctx, id)
	return user, err
}

func (r *entityResolver) FindUserByOAuthAccessToken(ctx context.Context, oAuthAccessToken string) (*model.User, error) {
	user, err := r.Resolver.Repository.GetUserByAuthToken(ctx, oAuthAccessToken)
	return user, err
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
