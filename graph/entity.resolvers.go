package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_users/graph/generated"
	"github.com/KnightHacks/knighthacks_users/graph/model"
)

// FindHackathonApplicationByID is the resolver for the findHackathonApplicationByID field.
func (r *entityResolver) FindHackathonApplicationByID(ctx context.Context, id string) (*model.HackathonApplication, error) {
	return &model.HackathonApplication{ID: id}, nil
}

// FindUserByID is the resolver for the findUserByID field.
func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := r.Resolver.Repository.GetUserByID(ctx, id)
	return user, err
}

// FindUserByOAuthUIDAndOAuthProvider is the resolver for the findUserByOAuthUIDAndOAuthProvider field.
func (r *entityResolver) FindUserByOAuthUIDAndOAuthProvider(ctx context.Context, oAuthUID string, oAuthProvider models.Provider) (*model.User, error) {
	user, err := r.Resolver.Repository.GetUserByOAuthUID(ctx, oAuthUID, oAuthProvider)
	return user, err
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
