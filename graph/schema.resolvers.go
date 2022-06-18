package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_users/graph/generated"
	"github.com/KnightHacks/knighthacks_users/graph/model"
)

func (r *mutationResolver) Register(ctx context.Context, provider models.Provider, encryptedOauthAccessToken string, input model.NewUser) (*model.RegistrationPayload, error) {
	// Decode the encrypted OAuth AccessToken from base64
	b, err := base64.URLEncoding.DecodeString(encryptedOauthAccessToken)
	if err != nil {
		return nil, err
	}
	// Decrypt the decoded access token using AES-256 decryption
	accessToken, err := r.Auth.DecryptAccessToken(string(b))
	if err != nil {
		return nil, err
	}
	// Using the access token retrieve the OAuth provided UID of the user
	uid, err := r.Auth.GetUID(ctx, provider, string(accessToken))
	if err != nil {
		return nil, err
	}
	// Create the user using the UID to check against duplicate accounts
	user, err := r.Repository.CreateUser(ctx, &model.OAuth{UID: uid, Provider: provider}, &input)
	if err != nil {
		// TODO: Possibly do some error handling hear to filter sql errors out
		return nil, err
	}

	refresh, access, err := r.Auth.NewTokens(user.ID, user.Role.String())
	if err != nil {
		return nil, err
	}
	payload := &model.RegistrationPayload{
		User:         user,
		RefreshToken: refresh,
		AccessToken:  access,
	}
	return payload, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAuthRedirectLink(ctx context.Context, provider models.Provider) (string, error) {
	return r.Auth.GetAuthCodeURL(provider), nil
}

func (r *queryResolver) Login(ctx context.Context, provider models.Provider, code string) (*model.LoginPayload, error) {
	// Using the OAuth code provided exchange the code for an access token
	token, err := r.Auth.ExchangeCode(ctx, provider, code)
	if err != nil {
		return nil, err
	}
	if !token.Valid() {
		// this shouldn't happen unless there was man-in-the-middle tampering to the HTTP request involved
		return nil, errors.New("auth token not valid, nice try hacker")
	}
	log.Printf("accessToken=%s, refreshToken=%s, type=%s, expiry=%s\n", token.AccessToken, token.RefreshToken, token.Type(), token.Expiry)
	// Get the user by their OAuth ID, if the user == nil then the user hasn't created an account yet, but will using the Register function
	uid, err := r.Auth.GetUID(ctx, provider, token.AccessToken)
	if err != nil {
		return nil, err
	}
	user, err := r.Repository.GetUserByOAuthUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	payload := model.LoginPayload{}
	if user != nil {
		// Set the user since they exist
		payload.User = user
		payload.AccountExists = true

		refresh, access, err := r.Auth.NewTokens(user.ID, user.Role.String())
		if err != nil {
			return nil, err
		}
		payload.RefreshToken = &refresh
		payload.AccessToken = &access
	} else {
		// Using AES-256 encryption, encrypt the access token to protect against packet sniffing
		encryptAccessTokenBytes := r.Auth.EncryptAccessToken(token.AccessToken)
		log.Printf("bytes=%v\n", encryptAccessTokenBytes)

		// Using base64 encoding, encode the access token to be able to be sent using alphanumeric character over HTTP
		encodedAccessToken := base64.URLEncoding.EncodeToString(encryptAccessTokenBytes)
		log.Printf("string=%v\n", encodedAccessToken)

		payload.EncryptedOAuthAccessToken = &encodedAccessToken
	}

	// The idea behind the last if statement is to return the user if it exists,
	// if the user does not exist then we should supply the user with the encodedAccessToken
	// to be able to Register an account
	return &payload, nil
}

func (r *queryResolver) RefreshJwt(ctx context.Context, refreshToken string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	return r.Repository.GetUserByID(ctx, id)
}

func (r *queryResolver) SearchUser(ctx context.Context, name string) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) FullName(ctx context.Context, obj *model.User) (string, error) {
	return fmt.Sprintf("%s %s", obj.FirstName, obj.LastName), nil
}

func (r *userResolver) OAuth(ctx context.Context, obj *model.User) (*model.OAuth, error) {
	return r.Repository.GetOAuth(ctx, obj.ID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
