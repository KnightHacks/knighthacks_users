package database

import (
	"context"
	"fmt"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewDatabaseRepository(ctx context.Context, databasePool *pgxpool.Pool) (*DatabaseRepository, error) {
	if databasePool == nil {
		return nil, fmt.Errorf("cannot create DatabaseRepository with nil databasePool")
	}
	databaseRepository := &DatabaseRepository{
		DatabasePool:      databasePool,
		PronounMap:        map[int]model.Pronouns{},
		PronounReverseMap: map[model.Pronouns]int{},
	}
	if err := databaseRepository.DatabasePool.Ping(ctx); err != nil {
		return nil, err
	}
	if err := databaseRepository.LoadPronouns(ctx); err != nil {
		return nil, err
	}
	return databaseRepository, nil
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
