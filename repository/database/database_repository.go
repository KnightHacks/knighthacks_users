package database

import (
	"context"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
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
