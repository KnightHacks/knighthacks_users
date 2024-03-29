package database

import (
	"context"
	"errors"
	"github.com/KnightHacks/knighthacks_shared/database"
	sharedModels "github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/jackc/pgx/v5"
	"strconv"
)

// CreateUser Creates a user in the database and returns the new user struct
//
// The NewUser input struct contains all nillable fields so the following function
// must be able to run regardless of whether of it's input, that is why there is a
// lot of pointers for nil safety purposes
func (r *DatabaseRepository) CreateUser(ctx context.Context, oAuth *model.OAuth, input *model.NewUser) (*model.User, error) {
	var pronouns *model.Pronouns = nil
	if input.Pronouns != nil {
		// input.Pronouns is a PronounsInput struct which a direct copy of the Pronouns struct, so we need to copy its fields
		pronouns = &model.Pronouns{
			Subjective: input.Pronouns.Subjective,
			Objective:  input.Pronouns.Objective,
		}
	}

	user := &model.User{
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Email:             input.Email,
		PhoneNumber:       input.PhoneNumber,
		Pronouns:          pronouns,
		Age:               input.Age,
		Role:              sharedModels.RoleNormal,
		OAuth:             oAuth,
		Race:              input.Race,
		YearsOfExperience: input.YearsOfExperience,
		ShirtSize:         input.ShirtSize,
		Gender:            input.Gender,
	}

	// Begins the database transaction
	err := pgx.BeginTxFunc(ctx, r.DatabasePool, pgx.TxOptions{}, func(tx pgx.Tx) error {
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
		var pronounIdPtr *int
		if pronouns != nil {
			pronounIdPtr, err = r.GetOrCreatePronoun(ctx, tx, *pronouns, input)
			if err != nil {
				return err
			}
		}
		// Insert new user into database
		userIdInt, err := r.InsertUser(ctx, tx, input, pronounIdPtr, oAuth)
		if err != nil {
			return err
		}

		// Insert MLH Terms
		if input.Mlh != nil {
			if err = r.InsertMLHTerms(ctx, tx, userIdInt, input.Mlh); err != nil {
				return err
			}
			user.Mlh = &model.MLHTerms{
				SendMessages:  input.Mlh.SendMessages,
				CodeOfConduct: input.Mlh.CodeOfConduct,
				ShareInfo:     input.Mlh.ShareInfo,
			}
		}

		// Insert Education Info
		if input.EducationInfo != nil {
			if err = r.InsertEducationInfo(ctx, tx, userIdInt, input.EducationInfo); err != nil {
				return err
			}
			user.EducationInfo = &model.EducationInfo{
				Name:           input.EducationInfo.Name,
				GraduationDate: input.EducationInfo.GraduationDate,
				Major:          input.EducationInfo.Major,
				Level:          input.EducationInfo.Level,
			}
		}

		// Insert Mailing Address Data
		if input.MailingAddress != nil {
			if err = r.InsertMailingAddress(ctx, tx, userIdInt, input.MailingAddress); err != nil {
				return err
			}
			user.MailingAddress = &model.MailingAddress{
				Country:      input.MailingAddress.Country,
				State:        input.MailingAddress.State,
				City:         input.MailingAddress.City,
				PostalCode:   input.MailingAddress.PostalCode,
				AddressLines: input.MailingAddress.AddressLines,
			}
		}

		user.ID = strconv.Itoa(userIdInt)
		return nil
	})
	if err != nil {
		return nil, err
	}
	// TODO: look into the case where the userId is not scanned in
	return user, nil
}

func (r *DatabaseRepository) InsertUser(ctx context.Context, queryable database.Queryable, input *model.NewUser, pronounIdPtr *int, oAuth *model.OAuth) (int, error) {
	// TODO: Possibly change ID type to int to stop this hacky fix?
	// insert user into database and return their ID

	var raceStringArray []string
	if input.Race != nil && len(input.Race) > 0 {
		raceStringArray = make([]string, 0, len(input.Race))
		for _, race := range input.Race {
			raceStringArray = append(raceStringArray, race.String())
		}
	}
	var userIdInt int
	err := queryable.QueryRow(ctx, "INSERT INTO users (first_name, last_name, email, phone_number, age, pronoun_id, oauth_uid, oauth_provider, role, years_of_experience, shirt_size, race, gender) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id",
		input.FirstName,
		input.LastName,
		input.Email,
		input.PhoneNumber,
		input.Age,
		pronounIdPtr,
		oAuth.UID,
		oAuth.Provider.String(),
		sharedModels.RoleNormal,
		input.YearsOfExperience,
		input.ShirtSize,
		raceStringArray,
		input.Gender,
	).Scan(&userIdInt)
	return userIdInt, err
}

func (r *DatabaseRepository) InsertMLHTerms(ctx context.Context, queryable database.Queryable, userId int, input *model.MLHTermsInput) error {
	_, err := queryable.Exec(ctx, "INSERT INTO mlh_terms (user_id, send_messages, share_info, code_of_conduct) VALUES ($1, $2, $3, $4)",
		userId,
		input.SendMessages,
		input.ShareInfo,
		input.CodeOfConduct,
	)
	return err
}

func (r *DatabaseRepository) InsertEducationInfo(ctx context.Context, queryable database.Queryable, userId int, input *model.EducationInfoInput) error {
	_, err := queryable.Exec(ctx, "INSERT INTO education_info (user_id, name, major, graduation_date, level) VALUES ($1, $2, $3, $4, $5)",
		userId,
		input.Name,
		input.Major,
		input.GraduationDate.UTC(),
		input.Level,
	)
	return err
}

func (r *DatabaseRepository) InsertMailingAddress(ctx context.Context, queryable database.Queryable, userId int, input *model.MailingAddressInput) error {
	_, err := queryable.Exec(ctx, "INSERT INTO mailing_addresses (user_id, country, state, city, postal_code, address_lines) VALUES ($1, $2, $3, $4, $5, $6)",
		userId,
		input.Country,
		input.State,
		input.City,
		input.PostalCode,
		input.AddressLines,
	)
	return err
}
