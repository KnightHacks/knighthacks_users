package database

import (
	"context"
	"errors"
	"github.com/KnightHacks/knighthacks_shared/database"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/jackc/pgx/v5"
)

/*
$$$$$$$\
$$  __$$\
$$ |  $$ | $$$$$$\   $$$$$$\  $$$$$$$\   $$$$$$\  $$\   $$\ $$$$$$$\   $$$$$$$\
$$$$$$$  |$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\ $$ |  $$ |$$  __$$\ $$  _____|
$$  ____/ $$ |  \__|$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |$$ |  $$ |\$$$$$$\
$$ |      $$ |      $$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ | \____$$\
$$ |      $$ |      \$$$$$$  |$$ |  $$ |\$$$$$$  |\$$$$$$  |$$ |  $$ |$$$$$$$  |
\__|      \__|       \______/ \__|  \__| \______/  \______/ \__|  \__|\_______/
*/

// GetByPronouns gets the sql row id for the pronouns associated with the input
func (r *DatabaseRepository) GetByPronouns(pronouns model.Pronouns) (int, bool) {
	id, exist := r.PronounReverseMap[pronouns]
	return id, exist
}

// GetById gets the pronouns by the sql row id
func (r *DatabaseRepository) GetById(id int) (model.Pronouns, bool) {
	pronouns, exist := r.PronounMap[id]
	return pronouns, exist
}

// Set inputs the pronouns into the bidirectional map
func (r *DatabaseRepository) Set(id int, pronouns model.Pronouns) {
	r.PronounMap[id] = pronouns
	r.PronounReverseMap[pronouns] = id
}

func (r *DatabaseRepository) GetPronouns(ctx context.Context, queryable database.Queryable, pronounId int) error {
	pronouns, exists := r.GetById(pronounId)
	// does the pronoun not exist in the local cache?
	if !exists {
		// retrieve the pronoun from the database
		err := queryable.QueryRow(ctx, "SELECT subjective, objective FROM pronouns WHERE id = $1", pronounId).Scan(
			&pronouns.Subjective,
			&pronouns.Objective,
		)
		if err != nil {
			return err
		}
		// set the pronoun in the local cache
		r.Set(pronounId, pronouns)
	}
	return nil
}

func (r *DatabaseRepository) GetOrCreatePronoun(ctx context.Context, queryable database.Queryable, pronouns model.Pronouns, input *model.NewUser) (*int, error) {
	pronounId, exists := r.GetByPronouns(pronouns)
	// if the pronoun does not exist in the local cache
	if !exists {
		// check if the pronoun exists in the database
		err := queryable.QueryRow(ctx, "SELECT id FROM pronouns WHERE subjective=$1 AND objective=$2",
			input.Pronouns.Subjective,
			input.Pronouns.Objective,
		).Scan(&pronounId)

		pronounExist := true
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// pronoun does not exist in the database
				pronounExist = false
			} else {
				return nil, err
			}
		}

		if !pronounExist {
			// since the new pronoun does not exist in the database, we insert it
			err = queryable.QueryRow(ctx, "INSERT INTO pronouns (subjective, objective) VALUES ($1, $2) RETURNING id",
				input.Pronouns.Subjective,
				input.Pronouns.Objective,
			).Scan(&pronounId)
		}

		if err != nil {
			return nil, err
		}
		// set the pronoun cache
		r.Set(pronounId, pronouns)
	}
	return &pronounId, nil
}

func (r *DatabaseRepository) LoadPronouns(ctx context.Context) error {
	rows, err := r.DatabasePool.Query(ctx, "SELECT id, subjective, objective FROM pronouns")
	if err != nil {
		return err
	}

	for rows.Next() {
		var pronouns model.Pronouns
		var id int
		err = rows.Scan(&id, &pronouns.Subjective, &pronouns.Objective)
		if err != nil {
			return err
		}
		r.Set(id, pronouns)
	}

	return nil
}
