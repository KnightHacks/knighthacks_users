package database

import (
	"context"
	"errors"
	sharedModels "github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/jackc/pgx/v4"
)

/*
 $$$$$$\             $$\           $$\   $$\
$$  __$$\            $$ |          $$ |  $$ |
$$ /  \__| $$$$$$\ $$$$$$\         $$ |  $$ | $$$$$$$\  $$$$$$\   $$$$$$\
$$ |$$$$\ $$  __$$\\_$$  _|        $$ |  $$ |$$  _____|$$  __$$\ $$  __$$\
$$ |\_$$ |$$$$$$$$ | $$ |          $$ |  $$ |\$$$$$$\  $$$$$$$$ |$$ |  \__|
$$ |  $$ |$$   ____| $$ |$$\       $$ |  $$ | \____$$\ $$   ____|$$ |
\$$$$$$  |\$$$$$$$\  \$$$$  |      \$$$$$$  |$$$$$$$  |\$$$$$$$\ $$ |
 \______/  \_______|  \____/        \______/ \_______/  \_______|\__|
*/

func (r *DatabaseRepository) GetUsers(ctx context.Context, first int, after string) ([]*model.User, int, error) {
	users := make([]*model.User, 0, first)
	var totalCount int
	err := r.DatabasePool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		rows, err := tx.Query(
			ctx,
			"SELECT id, first_name, last_name, email, phone_number, pronoun_id, age, role FROM users WHERE id > $1 ORDER BY `id` DESC LIMIT $2",
			after,
			first,
		)
		if err != nil {
			return err
		}
		for rows.Next() {
			var user model.User

			pronounId, err := ScanUser(&user, rows)
			if err != nil {
				return err
			}
			// user has pronouns, but we don't know what they are
			if pronounId != nil {
				err = r.getPronouns(ctx, tx, *pronounId)
				if err != nil {
					return err
				}
			}
			users = append(users, &user)
		}

		if err = rows.Err(); err != nil {
			return err
		}
		row := tx.QueryRow(ctx, "SELECT COUNT(*) FROM users")
		if err != nil {
			return err
		}

		return row.Scan(&totalCount)
	})
	return users, totalCount, err
}

// GetUserByID returns the user by their id
func (r *DatabaseRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return r.getUser(
		ctx,
		`SELECT id, first_name, last_name, email, phone_number, pronoun_id, age, role FROM users WHERE id = $1 LIMIT 1`,
		id,
	)
}

// GetUserByOAuthUID returns the user by their oauth auth token
func (r *DatabaseRepository) GetUserByOAuthUID(ctx context.Context, oAuthUID string, provider sharedModels.Provider) (*model.User, error) {
	return r.getUser(
		ctx,
		`SELECT id, first_name, last_name, email, phone_number, pronoun_id, age, role FROM users WHERE oauth_uid=cast($1 as varchar) AND oauth_provider=$2 LIMIT 1`,
		oAuthUID,
		provider,
	)
}

func (r *DatabaseRepository) getUserWithTx(ctx context.Context, query string, tx pgx.Tx, args ...interface{}) (*model.User, error) {
	var user model.User

	if tx == nil {
		var err error
		tx, err = r.DatabasePool.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Commit(ctx)
	}

	pronounId, err := ScanUser(&user, tx.QueryRow(ctx, query, args...))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.UserNotFound
		}
		return nil, err
	}

	// if the user has their pronouns set
	if pronounId != nil {
		if err = r.getPronouns(ctx, tx, *pronounId); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (r *DatabaseRepository) getUser(ctx context.Context, query string, args ...interface{}) (*model.User, error) {
	return r.getUserWithTx(ctx, query, nil, args...)
}

func (r *DatabaseRepository) SearchUser(ctx context.Context, name string) ([]*model.User, error) {
	const limit int = 10
	users := make([]*model.User, 0, limit)

	err := r.DatabasePool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, "SELECT * from users WHERE to_tsvector(first_name || ' ' || last_name) @@ to_tsquery('$1:*') LIMIT $2", name, limit)
		if err != nil {
			return err
		}
		for rows.Next() {
			var user model.User

			pronounId, err := ScanUser(&user, rows)
			if err != nil {
				return err
			}
			if pronounId != nil {
				err = r.getPronouns(ctx, tx, *pronounId)
			}
			if err != nil {
				return err
			}
			users = append(users, &user)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetOAuth returns the model.OAuth object that is associated with the user's id
// Used by the OAuth force resolver, this is not a common operation so making this
// a force resolver is a good idea
func (r *DatabaseRepository) GetOAuth(ctx context.Context, id string) (*model.OAuth, error) {
	var oAuth model.OAuth
	err := r.DatabasePool.QueryRow(ctx, "SELECT oauth_uid, oauth_provider FROM users WHERE id = $1", id).Scan(&oAuth.UID, &oAuth.Provider)
	if err != nil {
		return nil, err
	}
	return &oAuth, err
}

func (r *DatabaseRepository) GetUserMailingAddress(ctx context.Context, userId string) (*model.MailingAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (r *DatabaseRepository) GetUserMLHTerms(ctx context.Context, userId string) (*model.MLHTerms, error) {
	//TODO implement me
	panic("implement me")
}
