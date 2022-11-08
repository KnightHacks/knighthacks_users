package database

import (
	"context"
	"errors"
	"github.com/KnightHacks/knighthacks_shared/database"
	sharedModels "github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
)

// DatabaseRepository
// Implements the Repository interface's functions
//
// PronounMap & PronounReverseMap are the 2 maps that implement a bidirectional
// map to handle cached pronouns in the database to remove the need to do a SQL join
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

// CreateUser Creates a user in the database and returns the new user struct
//
// The NewUser input struct contains all nillable fields so the following function
// must be able to run regardless of whether of it's input, that is why there is a
// lot of pointers for nil safety purposes
func (r *DatabaseRepository) CreateUser(ctx context.Context, oAuth *model.OAuth, input *model.NewUser) (*model.User, error) {
	var userId string
	var pronouns model.Pronouns
	if input.Pronouns != nil {
		// input.Pronouns is a PronounsInput struct which a direct copy of the Pronouns struct, so we need to copy its fields
		pronouns = model.Pronouns{
			Subjective: input.Pronouns.Subjective,
			Objective:  input.Pronouns.Objective,
		}
	}

	// Begins the database transaction
	err := r.DatabasePool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		// Detects whether the user with the oauth_uid, for GitHub that is their github ID already exists, if
		// the use already exists we return an UserAlreadyExists error
		var discoveredId = new(int)
		err := tx.QueryRow(ctx, "SELECT id FROM users WHERE oauth_uid=$1 AND oauth_provider=$2 LIMIT 1", oAuth.UID, oAuth.Provider.String()).Scan(discoveredId)
		if err == nil && discoveredId != nil {
			return repository.UserAlreadyExists
		}
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		// Get pronouns
		pronounIdPtr, err := r.GetOrCreatePronoun(ctx, tx, pronouns, input)
		if err != nil {
			return err
		}

		userIdInt, err := r.InsertUser(ctx, tx, input, pronounIdPtr, oAuth)
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
		Pronouns:    &pronouns,
		Age:         input.Age,
		OAuth:       oAuth,
	}, nil
}

func (r *DatabaseRepository) InsertUser(ctx context.Context, queryable database.Queryable, input *model.NewUser, pronounIdPtr *int, oAuth *model.OAuth) (int, error) {
	// TODO: Possibly change ID type to int to stop this hacky fix?
	// insert user into database and return their ID
	var userIdInt int
	err := queryable.QueryRow(ctx, "INSERT INTO users (first_name, last_name, email, phone_number, age, pronoun_id, oauth_uid, oauth_provider, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		input.FirstName,
		input.LastName,
		input.Email,
		input.PhoneNumber,
		input.Age,
		pronounIdPtr,
		oAuth.UID,
		oAuth.Provider.String(),
		sharedModels.RoleNormal,
	).Scan(&userIdInt)
	return userIdInt, err
}

func (r *DatabaseRepository) DeleteUser(ctx context.Context, id string) (bool, error) {

	//query the row using the id with a DELETE statment
	commandTag, err := r.DatabasePool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)

	//err should return a nil value, if not throw error
	if err != nil {
		return false, err
	}

	//there should be one row affected, if not throw error
	if commandTag.RowsAffected() != 1 {
		return false, repository.UserNotFound
	}

	// then no error
	return true, nil
}

type Scannable interface {
	Scan(dest ...interface{}) error
}

func ScanUser[T Scannable](user *model.User, scannable T) (*int, error) {
	var pronounId *int32
	var userIdInt int
	err := scannable.Scan(
		&userIdInt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNumber,
		&pronounId,
		&user.Age,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}
	user.ID = strconv.Itoa(userIdInt)
	if pronounId != nil {
		i := int(*pronounId)
		return &i, nil
	}
	return nil, nil
}
