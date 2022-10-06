package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/KnightHacks/knighthacks_shared/auth"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_shared/pagination"
	"github.com/KnightHacks/knighthacks_shared/utils"
	"github.com/KnightHacks/knighthacks_users/graph/generated"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"net/http"
	"net/url"
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

	refresh, access, err := r.Auth.NewTokens(user.ID, user.Role)
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

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdatedUser) (*model.User, error) {
	if input.FirstName == nil && input.LastName == nil && input.Email == nil && input.PhoneNumber == nil && input.Pronouns == nil && input.Age == nil {
		return nil, fmt.Errorf("no field has been updated")
	}

	claims, ok := ctx.Value("AuthorizationUserClaims").(*auth.UserClaims)
	if !ok {
		return nil, errors.New("unable to retrieve user claims, most likely forgot to set @hasRole directive")
	}
	if claims.Role != models.RoleAdmin && claims.Id != id {
		return nil, errors.New("unauthorized to update user that is not you")
	}

	return r.Repository.UpdateUser(ctx, id, &input)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	claims, ok := ctx.Value("AuthorizationUserClaims").(*auth.UserClaims)
	if !ok {
		return false, errors.New("unable to retrieve user claims, most likely forgot to set @hasRole directive")
	}
	if claims.Role != models.RoleAdmin && claims.Id != id {
		return false, errors.New("unauthorized to update user that is not you")
	}
	return r.Repository.DeleteUser(ctx, id)
}

func (r *queryResolver) GetAuthRedirectLink(ctx context.Context, provider models.Provider) (string, error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return "", err
	}

	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)

	// TODO: check into enabling secure behind proxy in production
	ginContext.SetSameSite(http.SameSiteNoneMode)
	ginContext.SetCookie("oauthstate", state, 60*10, "/", "", false, true)
	ginContext.Header("Access-Control-Allow-Credentials", "true")
	return r.Auth.GetAuthCodeURL(provider, state), nil
}

func (r *queryResolver) Login(ctx context.Context, provider models.Provider, code string, state string) (*model.LoginPayload, error) {
	// todo: this should probably be cleaned up, been at this shit for hours, please god.. no more
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	cookieHeader := ginContext.GetHeader("oauthstate")
	cookieHeader, err = url.QueryUnescape(cookieHeader)
	if err != nil {
		return nil, err
	}

	if cookieHeader != state {
		return nil, fmt.Errorf("invalid oauth state")
	}

	// Using the OAuth code provided exchange the code for an access token
	token, err := r.Auth.ExchangeCode(ctx, provider, code)
	if err != nil {
		return nil, err
	}
	if !token.Valid() {
		// this shouldn't happen unless there was man-in-the-middle tampering to the HTTP request involved
		return nil, errors.New("auth token not valid, nice try hacker")
	}
	// Get the user by their OAuth ID, if the user == nil then the user hasn't created an account yet, but will using the Register function
	uid, err := r.Auth.GetUID(ctx, provider, token.AccessToken)
	if err != nil {
		return nil, err
	}
	user, err := r.Repository.GetUserByOAuthUID(ctx, uid, provider)
	if err != nil && !errors.Is(err, repository.UserNotFound) {
		return nil, err
	}
	payload := model.LoginPayload{}
	if user != nil {
		// Set the user since they exist
		payload.User = user
		payload.AccountExists = true

		refresh, access, err := r.Auth.NewTokens(user.ID, user.Role)
		if err != nil {
			return nil, err
		}
		payload.RefreshToken = &refresh
		payload.AccessToken = &access
	} else {
		// Using AES-256 encryption, encrypt the access token to protect against packet sniffing
		encryptAccessTokenBytes := r.Auth.EncryptAccessToken(token.AccessToken)

		// Using base64 encoding, encode the access token to be able to be sent using alphanumeric character over HTTP
		encodedAccessToken := base64.URLEncoding.EncodeToString(encryptAccessTokenBytes)

		payload.EncryptedOAuthAccessToken = &encodedAccessToken
	}

	// The idea behind the last if statement is to return the user if it exists,
	// if the user does not exist then we should supply the user with the encodedAccessToken
	// to be able to Register an account
	return &payload, nil
}

func (r *queryResolver) RefreshJwt(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenUserClaims, err := r.Auth.ParseJWT(refreshToken, auth.RefreshTokenType)
	if err != nil {
		// TODO: special handler for auth.TokenNotValid error
		// if the err is auth.TokenNotValid then the user must login again
		return "", err
	}
	token, err := r.Auth.NewAccessToken(refreshTokenUserClaims.UserID, refreshTokenUserClaims.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) Users(ctx context.Context, first int, after *string) (*model.UsersConnection, error) {
	a, err := pagination.DecodeCursor(after)
	if err != nil {
		return nil, err
	}
	users, total, err := r.Repository.GetUsers(ctx, first, a)
	if err != nil {
		return nil, err
	}

	return &model.UsersConnection{
		TotalCount: total,
		PageInfo:   pagination.GetPageInfo(users[0].ID, users[len(users)-1].ID),
		Users:      users,
	}, nil
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	return r.Repository.GetUserByID(ctx, id)
}

func (r *queryResolver) SearchUser(ctx context.Context, name string) ([]*model.User, error) {
	if !utils.IsASCII(name) {
		// TODO: how to handle non ascii names? do they exist? idk
		return nil, fmt.Errorf("the name must include only ascii characters")
	}

	return r.Repository.SearchUser(ctx, name)
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userClaims, err := auth.UserClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return r.Entity().FindUserByID(ctx, userClaims.UserID)
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
