package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/KnightHacks/knighthacks_shared/database"
	"github.com/KnightHacks/knighthacks_shared/utils"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/jackc/pgx/v5"
	"strconv"
	"time"
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

// UpdatePronouns(ctx context.Context, id string, pronoun *model.PronounsInput, tx pgx.Tx) error
// func (ctx context.Context, id string, t *T, tx pgx.Tx)
type UpdateFunc[T any] func(ctx context.Context, id string, input T, tx pgx.Tx) error

// Validate Yes, I understand this function is generic hell but there is no other way to do it.
// *any does not work bc when you use *any it passes your generic type into any and not as *any.
// that would make it a double pointer and not a single pointer.
func Validate[T *string |
	*float64 |
	*model.ShirtSize |
	*int |
	[]*string |
	*model.PronounsInput |
	*model.MailingAddressUpdate |
	*model.EducationInfoUpdate |
	*model.MLHTermsUpdate |
	[]model.Race](ctx context.Context, tx pgx.Tx, id string, input T, updateFunc UpdateFunc[T]) error {
	if input != nil {
		err := updateFunc(ctx, id, input, tx)
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
	err = pgx.BeginFunc(ctx, r.DatabasePool, func(tx pgx.Tx) error {
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
		if err = Validate(ctx, tx, id, input.Pronouns, r.UpdatePronouns); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.Age, r.UpdateAge); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.EducationInfo, r.UpdateEducationInfo); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.Mlh, r.UpdateMLHTerms); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.MailingAddress, r.UpdateMailingAddress); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.ShirtSize, r.UpdateShirtSize); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.Gender, r.UpdateGender); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.Race, r.UpdateRace); err != nil {
			return err
		}
		if err = Validate(ctx, tx, id, input.YearsOfExperience, r.UpdateYearsOfExperience); err != nil {
			return err
		}

		user, err = r.GetUserWithQueryable(ctx,
			`SELECT id, first_name, last_name, email, phone_number, pronoun_id, age, role, gender, race, shirt_size, years_of_experience FROM users WHERE id = $1 LIMIT 1`,
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
func (r *DatabaseRepository) UpdateFirstName(ctx context.Context, id string, first *string, tx pgx.Tx) error {

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
func (r *DatabaseRepository) UpdateLastName(ctx context.Context, id string, last *string, tx pgx.Tx) error {
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
func (r *DatabaseRepository) UpdateEmail(ctx context.Context, id string, email *string, tx pgx.Tx) error {
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
func (r *DatabaseRepository) UpdatePhoneNumber(ctx context.Context, id string, number *string, tx pgx.Tx) error {
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

// UpdateShirtSize updates shirt size
func (r *DatabaseRepository) UpdateShirtSize(ctx context.Context, id string, shirtSize *model.ShirtSize, tx pgx.Tx) error {
	commandTag, err := tx.Exec(ctx, "UPDATE users SET shirt_size = $1 WHERE id = $2", shirtSize.String(), id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdateYearsOfExperience updates years of experience
func (r *DatabaseRepository) UpdateYearsOfExperience(ctx context.Context, id string, years *float64, tx pgx.Tx) error {
	commandTag, err := tx.Exec(ctx, "UPDATE users SET years_of_experience = $1 WHERE id = $2", *years, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdateRace updates race
func (r *DatabaseRepository) UpdateRace(ctx context.Context, id string, races []model.Race, tx pgx.Tx) error {
	var raceStringArray []string
	if races != nil && len(races) > 0 {
		raceStringArray = make([]string, 0, len(races))
		for _, race := range races {
			raceStringArray = append(raceStringArray, race.String())
		}
	}
	commandTag, err := tx.Exec(ctx, "UPDATE users SET race = $1 WHERE id = $2", raceStringArray, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdateGender updates gender
func (r *DatabaseRepository) UpdateGender(ctx context.Context, id string, gender *string, tx pgx.Tx) error {
	commandTag, err := tx.Exec(ctx, "UPDATE users SET gender = $1 WHERE id = $2", *gender, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdateMLHTerms updates user's MLH Terms
func (r *DatabaseRepository) UpdateMLHTerms(ctx context.Context, id string, input *model.MLHTermsUpdate, tx pgx.Tx) error {
	var keys []any
	var values []any

	if input.CodeOfConduct != nil {
		keys = append(keys, "code_of_conduct")
		values = append(values, *input.CodeOfConduct)
	}
	if input.ShareInfo != nil {
		keys = append(keys, "share_info")
		values = append(values, *input.ShareInfo)
	}
	if input.SendMessages != nil {
		keys = append(keys, "send_messages")
		values = append(values, *input.SendMessages)
	}
	// Extra check to ensure there is something being sent, should never happen
	if len(keys) == 0 || len(values) == 0 || len(values) != len(keys) {
		return errors.New("something went wrong calculating keys and values for sql")
	}

	sql := fmt.Sprintf(`UPDATE mlh_terms SET %s WHERE user_id = $1`,
		database.GenerateUpdatePairs(keys, 2))

	combined := append(keys, values...)

	commandTag, err := tx.Exec(ctx, sql, id, combined)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdateEducationInfo updates user's Edu info
func (r *DatabaseRepository) UpdateEducationInfo(ctx context.Context, id string, input *model.EducationInfoUpdate, tx pgx.Tx) error {
	var keys []any
	var values []any

	if input.Level != nil {
		keys = append(keys, "level")
		values = append(values, *input.Level)
	}
	if input.Name != nil {
		keys = append(keys, "name")
		values = append(values, *input.Name)
	}
	if input.GraduationDate != nil {
		keys = append(keys, "graduation_date")
		values = append(values, *input.GraduationDate)
	}
	if input.Major != nil {
		keys = append(keys, "major")
		values = append(values, *input.Major)
	}
	// Extra check to ensure there is something being sent, should never happen
	if len(keys) == 0 || len(values) == 0 || len(values) != len(keys) {
		return errors.New("something went wrong calculating keys and values for sql")
	}

	sql := fmt.Sprintf(`UPDATE education_info SET %s WHERE user_id = $1`,
		database.GenerateUpdatePairs(keys, 2))

	combined := append([]any{id}, values...)

	commandTag, err := tx.Exec(ctx, sql, combined...)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return repository.UserNotFound
	}
	// then no error
	return nil
}

// UpdateMailingAddress updates user's mailing address
func (r *DatabaseRepository) UpdateMailingAddress(ctx context.Context, id string, input *model.MailingAddressUpdate, tx pgx.Tx) error {
	var keys []any
	var values []any

	if input.Country != nil {
		keys = append(keys, "country")
		values = append(values, *input.Country)
	}
	if input.State != nil {
		keys = append(keys, "state")
		values = append(values, *input.State)
	}
	if input.City != nil {
		keys = append(keys, "city")
		values = append(values, *input.City)
	}
	if input.PostalCode != nil {
		keys = append(keys, "postal_code")
		values = append(values, *input.PostalCode)
	}
	if input.AddressLines != nil && len(input.AddressLines) > 0 {
		keys = append(keys, "address_lines")
		values = append(values, input.AddressLines)
	}
	// Extra check to ensure there is something being sent, should never happen
	if len(keys) == 0 || len(values) == 0 || len(values) != len(keys) {
		return errors.New("something went wrong calculating keys and values for sql")
	}

	sql := fmt.Sprintf(`UPDATE mailing_addresses SET %s WHERE user_id = $1`,
		database.GenerateUpdatePairs(keys, 2))

	combined := append(keys, values...)

	commandTag, err := tx.Exec(ctx, sql, id, combined)
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

type Scannable interface {
	Scan(dest ...interface{}) error
}

func ScanUser[T Scannable](user *model.User, scannable T) (*int, error) {
	var pronounVal uint32
	pronounId := &pronounVal
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
		&user.Gender,
		&user.Race,
		&user.ShirtSize,
		&user.YearsOfExperience,
	)
	if err != nil {
		return nil, err
	}
	user.ID = strconv.Itoa(userIdInt)
	if pronounId == nil {
		return nil, nil
	}
	return utils.Ptr(int(*pronounId)), nil
}

func (r *DatabaseRepository) DeleteAPIKey(ctx context.Context, id string) error {
	_, err := r.DatabasePool.Exec(ctx, "DELETE FROM api_keys WHERE user_id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *DatabaseRepository) AddAPIKey(ctx context.Context, id string, key string) (*model.APIKey, error) {
	now := time.Now()
	_, err := r.DatabasePool.Exec(ctx, "INSERT INTO api_keys (user_id, key, created) VALUES ($1, $2, $3)", id, key, now)
	if err != nil {
		return nil, err
	}
	return &model.APIKey{Key: key, Created: time.Now()}, nil
}
