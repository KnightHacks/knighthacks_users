package database

import (
	"context"
	"fmt"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/jackc/pgx/v5"
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
	err := pgx.BeginFunc(ctx, r.DatabasePool, func(tx pgx.Tx) error {
		// Delete from hackathon_applications
		_, err := tx.Exec(ctx, "DELETE FROM hackathon_applications WHERE user_id = $1", id)
		if err != nil {
			return err
		}

		// Delete from mlh_terms
		_, err = tx.Exec(ctx, "DELETE FROM mlh_terms WHERE user_id = $1", id)
		if err != nil {
			return err
		}

		// Delete from meals
		_, err = tx.Exec(ctx, "DELETE FROM meals WHERE user_id = $1", id)
		if err != nil {
			return err
		}

		// Delete from mailing_addresses
		_, err = tx.Exec(ctx, "DELETE FROM mailing_addresses WHERE user_id = $1", id)
		if err != nil {
			return err
		}

		// Delete from hackathon_checkin
		_, err = tx.Exec(ctx, "DELETE FROM hackathon_checkin WHERE user_id = $1", id)
		if err != nil {
			return err
		}

		// Delete from hackathon_checkin
		_, err = tx.Exec(ctx, "DELETE FROM education_info WHERE user_id = $1", id)
		if err != nil {
			return err
		}

		// Delete from api_keys
		_, err = tx.Exec(ctx, "DELETE FROM api_keys WHERE user_id = $1", id)
		if err != nil {
			return err
		}

		// Delete from api_keys
		_, err = tx.Exec(ctx, "DELETE FROM event_attendance WHERE user_id = $1", id)
		if err != nil {
			return err
		}

		commandTag, err := r.DatabasePool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
		if err != nil {
			return err
		}

		//there should be one row affected, if not throw error
		if commandTag.RowsAffected() != 1 {
			return repository.UserNotFound
		}

		return nil
	})

	if err != nil {
		return false, err
	}
	// then no error
	return true, nil
}
