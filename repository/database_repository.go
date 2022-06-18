package repository

import (
	"context"
	"errors"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
)

var (
	UserNotFound      = errors.New("user not found")
	UserAlreadyExists = errors.New("user with id already exists")
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
func (r *DatabaseRepository) GetByPronouns(pronouns model.Pronouns) (int, bool) {
	id, exist := r.PronounReverseMap[pronouns]
	return id, exist
}

//GetById gets the pronouns by the sql row id
func (r *DatabaseRepository) GetById(id int) (model.Pronouns, bool) {
	pronouns, exist := r.PronounMap[id]
	return pronouns, exist
}

//Set inputs the pronouns into the bidirectional map
func (r *DatabaseRepository) Set(id int, pronouns model.Pronouns) {
	r.PronounMap[id] = pronouns
	r.PronounReverseMap[pronouns] = id
}

//GetUserByID returns the user by their id
func (r *DatabaseRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return r.getUser(ctx, "id", id)
}

// GetUserByOAuthUID returns the user by their oauth auth token
//TODO: possibly add Provider as argument?
func (r *DatabaseRepository) GetUserByOAuthUID(ctx context.Context, authToken string) (*model.User, error) {
	return r.getUser(ctx, "oauth_uid", authToken)
}

//GetOAuth returns the model.OAuth object that is associated with the user's id
//Used by the OAuth force resolver, this is not a common operation so making this
//a force resolver is a good idea
func (r *DatabaseRepository) GetOAuth(ctx context.Context, id string) (*model.OAuth, error) {
	var oAuth model.OAuth
	err := r.DatabasePool.QueryRow(ctx, "SELECT oauth_uid, oauth_provider FROM users WHERE id = $1", id).Scan(&oAuth.UID, &oAuth.Provider)
	if err != nil {
		return nil, err
	}
	return &oAuth, err
}

//getUser returns user by some column and value on the users table
func (r *DatabaseRepository) getUser(ctx context.Context, column string, value string) (*model.User, error) {
	var user model.User
	var pronounIdPtr *int
	err := r.DatabasePool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "SELECT first_name, last_name, email, phone_number, pronoun_id, age, role FROM users WHERE $1 = $2", column, value).Scan(
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PhoneNumber,
			pronounIdPtr,
			&user.Age,
			&user.Role,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return UserNotFound
			}
			return err
		}
		// if the user has their pronouns set
		if pronounIdPtr != nil {
			pronounId := *pronounIdPtr
			pronouns, exists := r.GetById(pronounId)
			// does the pronoun not exist in the local cache?
			if !exists {
				// retrieve the pronoun from the database
				err = tx.QueryRow(ctx, "SELECT subjective, objective FROM pronouns WHERE id = $1", pronounId).Scan(
					&pronouns.Subjective,
					&pronouns.Objective,
				)
				if err != nil {
					return err
				}
				// set the pronoun in the local cache
				r.Set(pronounId, pronouns)
			}
			user.Pronouns = &pronouns
		}
		return err
	})
	if err != nil {
		if errors.Is(err, UserNotFound) {
			// if the user does not exist then the user is nil, TODO: maybe return error?
			return nil, nil
		}
		return nil, err
	}
	return &user, err
}

//CreateUser Creates a user in the database and returns the new user struct
//
//The NewUser input struct contains all nillable fields so the following function
//must be able to run regardless of whether of it's input, that is why there is a
//lot of pointers for nil safety purposes
func (r *DatabaseRepository) CreateUser(ctx context.Context, oAuth *model.OAuth, input *model.NewUser) (*model.User, error) {
	var userId string
	var pronounsPtr *model.Pronouns
	if input.Pronouns != nil {
		// input.Pronouns is a PronounsInput struct which a direct copy of the Pronouns struct, so we need to copy its fields
		pronounsPtr = &model.Pronouns{
			Subjective: input.Pronouns.Subjective,
			Objective:  input.Pronouns.Objective,
		}
	}

	// Begins the database transaction
	err := r.DatabasePool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		// Detects whether or not the user with the oauth_uid, for github that is their github ID already exists, if
		// the use already exists we return an UserAlreadyExists error
		var discoveredId *int
		err := tx.QueryRow(ctx, "SELECT id FROM users WHERE oauth_uid=$1 AND oauth_provider=$2 LIMIT 1", oAuth.UID, oAuth.Provider.String()).Scan(discoveredId)
		if err == nil || discoveredId != nil {
			return UserAlreadyExists
		}
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		var pronounIdPtr *int

		// is pronoun being set?
		if pronounsPtr != nil {
			pronouns := *pronounsPtr
			pronounId, exists := r.GetByPronouns(pronouns)
			// if the pronoun does not exist in the local cache
			if !exists {
				// check if the pronoun exists in the database
				err = tx.QueryRow(ctx, "SELECT id FROM pronouns WHERE subjective=$1 AND objective=$2 RETURNING id",
					input.Pronouns.Subjective,
					input.Pronouns.Objective,
				).Scan(&pronounId)

				pronounExist := true
				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						// pronoun does not exist in the database
						pronounExist = false
					} else {
						return err
					}
				}
				if !pronounExist {
					// since the new pronoun does not exist in the database, we insert it
					err = tx.QueryRow(ctx, "INSERT INTO pronouns (subjective, objective) VALUES ($1, $2) RETURNING id",
						input.Pronouns.Subjective,
						input.Pronouns.Objective,
					).Scan(&pronounId)
				}

				if err != nil {
					return err
				}
				// set the pronoun cache
				r.Set(pronounId, pronouns)
			}

			pronounIdPtr = &pronounId
		}

		// TODO: Possibly change ID type to int to stop this hacky fix?
		// insert user into database and return their ID
		var userIdInt int
		err = tx.QueryRow(ctx, "INSERT INTO users (first_name, last_name, email, phone_number, age, pronoun_id, oauth_uid, oauth_provider, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
			input.FirstName,
			input.LastName,
			input.Email,
			input.PhoneNumber,
			input.Age,
			pronounIdPtr,
			oAuth.UID,
			oAuth.Provider.String(),
			model.RoleNormal,
		).Scan(&userIdInt)
		if err != nil {
			return err
		}
		userId = strconv.Itoa(userIdInt)
		return nil
	})
	if err != nil {
		return nil, err
	}
	// TODO: look into the case where the userId is not scanned in
	return &model.User{
		ID:          userId,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Pronouns:    pronounsPtr,
		Age:         input.Age,
		OAuth:       oAuth,
	}, nil
}
