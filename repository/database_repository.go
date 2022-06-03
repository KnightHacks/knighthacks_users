package repository

import (
	"context"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//DatabaseRepository
//Implements the Repository interface's functions
//
//PronounMap & PronounReverseMap are the 2 maps that implement a bidirectional
//map to handle cached pronouns in the database to remove the need to do a SQL join
type DatabaseRepository struct {
	DatabasePool      *pgxpool.Pool
	PronounMap        map[int]model.Pronouns
	PronounReverseMap map[model.Pronouns]int
}

func NewDatabaseRepository(databasePool *pgxpool.Pool) *DatabaseRepository {
	return &DatabaseRepository{
		DatabasePool:      databasePool,
		PronounMap:        map[int]model.Pronouns{},
		PronounReverseMap: map[model.Pronouns]int{},
	}
}

//GetByPronouns gets the sql row id for the pronouns associated with the input
func (r DatabaseRepository) GetByPronouns(pronouns model.Pronouns) (int, bool) {
	id, exist := r.PronounReverseMap[pronouns]
	return id, exist
}

//GetById gets the pronouns by the sql row id
func (r DatabaseRepository) GetById(id int) (model.Pronouns, bool) {
	pronouns, exist := r.PronounMap[id]
	return pronouns, exist
}

//Set inputs the pronouns into the bidirectional map
func (r *DatabaseRepository) Set(id int, pronouns model.Pronouns) {
	r.PronounMap[id] = pronouns
	r.PronounReverseMap[pronouns] = id
}

//GetUserByID returns the user by their id
func (r DatabaseRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return r.getUser(ctx, "id", id)
}

//GetUserByAuthToken returns the user by their oauth auth token
//TODO: possibly add Provider as argument?
func (r DatabaseRepository) GetUserByAuthToken(ctx context.Context, authToken string) (*model.User, error) {
	return r.getUser(ctx, "oauth_token", authToken)
}

//GetOAuth returns the model.OAuth object that is associated with the user's id
//Used by the OAuth force resolver, this is not a common operation so making this
//a force resolver is a good idea
func (r DatabaseRepository) GetOAuth(ctx context.Context, id string) (*model.OAuth, error) {
	var oAuth model.OAuth
	err := r.DatabasePool.QueryRow(ctx, "SELECT oauth_token, oauth_provider FROM users WHERE id = $1", id).Scan(&oAuth.AccessToken, &oAuth.Provider)
	if err != nil {
		return nil, err
	}
	return &oAuth, err
}

//getUser returns user by some key and value on the users table
func (r DatabaseRepository) getUser(ctx context.Context, key string, value string) (*model.User, error) {
	var user model.User
	var pronounId int
	err := r.DatabasePool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "SELECT first_name, last_name, email, phone_number, pronoun_id, age FROM users WHERE $1 = $2", key, value).Scan(
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PhoneNumber,
			&pronounId,
			&user.Age,
		)
		if err != nil {
			return err
		}
		pronouns, exists := r.GetById(pronounId)
		if !exists {
			err = tx.QueryRow(ctx, "SELECT subjective, objective FROM pronouns WHERE id = $1", pronounId).Scan(
				&pronouns.Subjective,
				&pronouns.Objective,
			)
			if err != nil {
				return err
			}
			r.Set(pronounId, pronouns)
		}
		user.Pronouns = &pronouns
		return err
	})
	if err != nil {
		return nil, err
	}
	return &user, err
}
