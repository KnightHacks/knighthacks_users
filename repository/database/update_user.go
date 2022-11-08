package database

import (
	"context"
	"errors"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/jackc/pgx/v4"
)

/*
$$\   $$\                 $$\            $$\                     $$\   $$\
$$ |  $$ |                $$ |           $$ |                    $$ |  $$ |
$$ |  $$ | $$$$$$\   $$$$$$$ | $$$$$$\ $$$$$$\    $$$$$$\        $$ |  $$ | $$$$$$$\  $$$$$$\   $$$$$$\
$$ |  $$ |$$  __$$\ $$  __$$ | \____$$\\_$$  _|  $$  __$$\       $$ |  $$ |$$  _____|$$  __$$\ $$  __$$\
$$ |  $$ |$$ /  $$ |$$ /  $$ | $$$$$$$ | $$ |    $$$$$$$$ |      $$ |  $$ |\$$$$$$\  $$$$$$$$ |$$ |  \__|
$$ |  $$ |$$ |  $$ |$$ |  $$ |$$  __$$ | $$ |$$\ $$   ____|      $$ |  $$ | \____$$\ $$   ____|$$ |
\$$$$$$  |$$$$$$$  |\$$$$$$$ |\$$$$$$$ | \$$$$  |\$$$$$$$\       \$$$$$$  |$$$$$$$  |\$$$$$$$\ $$ |
 \______/ $$  ____/  \_______| \_______|  \____/  \_______|       \______/ \_______/  \_______|\__|
          $$ |
          $$ |
          \__|
*/

type UpdateFunc[T any] func(ctx context.Context, id string, input T, tx pgx.Tx) error

func Validate[T any](ctx context.Context, tx pgx.Tx, id string, input *T, updateFunc UpdateFunc[T]) error {
	if input != nil {
		err := updateFunc(ctx, id, *input, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateUser
// update user add multiple parts go off of create user
// we will check whether the values in input are nil or empty strings, if not, we execute the update statement
func (r *DatabaseRepository) UpdateUser(ctx context.Context, id string, input *model.UpdatedUser) (*model.User, error) {
	var user *model.User
	var err error
	// checking to see if input is empty first
	if input.FirstName == nil && input.LastName == nil && input.Email == nil && input.PhoneNumber == nil && input.Pronouns == nil && input.Age == nil {
		return nil, errors.New("empty user field")
	}
	err = r.DatabasePool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if err = Validate(ctx, tx, id, input.FirstName, r.UpdateFirstName); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.LastName, r.UpdateLastName); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.Email, r.UpdateEmail); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.PhoneNumber, r.UpdatePhoneNumber); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, &input.Pronouns, r.UpdatePronouns); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, &input.Age, r.UpdateAge); err != nil {
			return err
		}

		user, err = r.getUserWithTx(ctx,
			`SELECT id, first_name, last_name, email, phone_number, pronoun_id, age, role FROM users WHERE id = $1 LIMIT 1`,
			tx,
			id)

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateFirstName this will update first name
func (r *DatabaseRepository) UpdateFirstName(ctx context.Context, id string, first string, tx pgx.Tx) error {

	commandTag, err := tx.Exec(ctx, "UPDATE users SET first_name = $1 WHERE id = $2", first, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	return nil
}

// UpdateLastName this function will update last name
func (r *DatabaseRepository) UpdateLastName(ctx context.Context, id string, last string, tx pgx.Tx) error {
	commandTag, err := tx.Exec(ctx, "UPDATE users SET last_name = $1 WHERE id = $2", last, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdateEmail updates email
func (r *DatabaseRepository) UpdateEmail(ctx context.Context, id string, email string, tx pgx.Tx) error {
	commandTag, err := tx.Exec(ctx, "UPDATE users SET email = $1 WHERE id = $2", email, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdatePhoneNumber updates user phone number
func (r *DatabaseRepository) UpdatePhoneNumber(ctx context.Context, id string, number string, tx pgx.Tx) error {
	commandTag, err := tx.Exec(ctx, "UPDATE users SET phone_number = $1 WHERE id = $2", number, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdateAge updates user age
func (r *DatabaseRepository) UpdateAge(ctx context.Context, id string, age *int, tx pgx.Tx) error {
	commandTag, err := tx.Exec(ctx, "UPDATE users SET age = $1 WHERE id = $2", *age, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdatePronouns updates user Pronouns
func (r *DatabaseRepository) UpdatePronouns(ctx context.Context, id string, pronoun *model.PronounsInput, tx pgx.Tx) error {
	// first find pronouns, if it doesn't exist this will add to database and then update user pronoun in database
	// copied from createUser
	var pronouns = model.Pronouns{
		Subjective: pronoun.Subjective,
		Objective:  pronoun.Objective,
	}
	pronounId, exists := r.GetByPronouns(pronouns)

	if !exists {
		// check if the pronoun exists in the database
		err := tx.QueryRow(ctx, "SELECT id FROM pronouns WHERE subjective=$1 AND objective=$2",
			pronoun.Subjective,
			pronoun.Objective,
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
				pronoun.Subjective,
				pronoun.Objective,
			).Scan(&pronounId)
		}

		if err != nil {
			return err
		}
		// set the pronoun cache
		r.Set(pronounId, pronouns)
	}

	commandTag, err := tx.Exec(ctx, "UPDATE users SET pronoun_id = $1 WHERE id = $2", pronounId, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}
