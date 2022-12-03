package integration_tests

import (
	"context"
	"flag"
	"fmt"
	shared_db_utils "github.com/KnightHacks/knighthacks_shared/database"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_shared/utils"
	model "github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

var integrationTest = flag.Bool("integration", false, "whether to run integration tests")
var databaseUri = flag.String("postgres-uri", "postgresql://postgres:test@localhost:5432/postgres", "postgres uri for running integration tests")

var databaseRepository *database.DatabaseRepository

type Test[A any, T any] struct {
	name    string
	args    A
	want    T
	wantErr bool
}

func TestMain(t *testing.M) {
	flag.Parse()
	// check if integration testing is disabled
	if *integrationTest == false {
		return
	}

	// connect to database
	var err error
	pool, err := shared_db_utils.ConnectWithRetries(*databaseUri)
	fmt.Printf("connecting to database, pool=%v, err=%v\n", pool, err)
	if err != nil {
		log.Fatalf("unable to connect to database err=%v\n", err)
	}

	databaseRepository, err = database.NewDatabaseRepository(context.Background(), pool)
	if err != nil {
		log.Fatalf("unable to initialize database repository err=%v\n", err)
	}
	os.Exit(t.Run())
}

func TestDatabaseRepository_AddAPIKey(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
		key string
	}
	tests := []Test[args, *model.APIKey]{
		{
			name: "add APIKey to Joe Bob",
			args: args{
				ctx: context.Background(),
				id:  "1",
				key: "12345abcdef",
			},
			want: &model.APIKey{
				Key: "12345abcdef",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey, err := databaseRepository.AddAPIKey(tt.args.ctx, tt.args.id, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(apiKey.Key, tt.want.Key) {
				t.Errorf("AddAPIKey() apiKey = %v, want %v", apiKey, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_CreateUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		oAuth *model.OAuth
		input *model.NewUser
	}
	tests := []Test[args, *model.User]{
		{
			name: "create thomas",
			args: args{
				ctx: context.Background(),
				oAuth: &model.OAuth{
					Provider: models.ProviderGithub,
					UID:      "100",
				},
				input: &model.NewUser{
					FirstName:   "Thomas",
					LastName:    "Bob",
					Email:       "thomas.bob@example.com",
					PhoneNumber: "100-203-9112",
					Pronouns: &model.PronounsInput{
						Subjective: "He",
						Objective:  "Him",
					},
					Age: utils.Ptr(21),
					MailingAddress: &model.MailingAddressInput{
						Country:    "United States",
						State:      "Florida",
						City:       "Orlando",
						PostalCode: "32765",
						AddressLines: []string{
							"1234 Joe Mama Row",
						},
					},
					Mlh: &model.MLHTermsInput{
						SendMessages:  true,
						CodeOfConduct: true,
						ShareInfo:     true,
					},
					ShirtSize:         model.ShirtSizeM,
					YearsOfExperience: utils.Ptr(3.5),
					EducationInfo: &model.EducationInfoInput{
						Name:           "University of Central Florida",
						GraduationDate: time.Date(2026, 12, 20, 0, 0, 0, 0, time.UTC),
						Major:          "Bachelors of Science",
						Level:          utils.Ptr(model.LevelOfStudyFreshman),
					},
					Gender: utils.Ptr("male"),
					Race:   []model.Race{model.RaceCaucasian, model.RaceAfricanAmerican},
				},
			},
			want: &model.User{
				//ID:          "", don't check for this
				FirstName:   "Thomas",
				LastName:    "Bob",
				Email:       "thomas.bob@example.com",
				PhoneNumber: "100-203-9112",
				Pronouns: &model.Pronouns{
					Subjective: "He",
					Objective:  "Him",
				},
				Age: utils.Ptr(21),
				MailingAddress: &model.MailingAddress{
					Country:    "United States",
					State:      "Florida",
					City:       "Orlando",
					PostalCode: "32765",
					AddressLines: []string{
						"1234 Joe Mama Row",
					},
				},
				Role:   models.RoleNormal,
				Gender: utils.Ptr("male"),
				Race:   []model.Race{model.RaceCaucasian, model.RaceAfricanAmerican},
				OAuth: &model.OAuth{
					Provider: models.ProviderGithub,
					UID:      "100",
				},
				Mlh: &model.MLHTerms{
					SendMessages:  true,
					CodeOfConduct: true,
					ShareInfo:     true,
				},
				ShirtSize:         model.ShirtSizeM,
				YearsOfExperience: utils.Ptr(3.5),
				EducationInfo: &model.EducationInfo{
					Name:           "University of Central Florida",
					GraduationDate: time.Date(2026, 12, 20, 0, 0, 0, 0, time.UTC),
					Major:          "Bachelors of Science",
					Level:          utils.Ptr(model.LevelOfStudyFreshman),
				},
				APIKey: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := databaseRepository.CreateUser(tt.args.ctx, tt.args.oAuth, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(user, tt.want) {
				t.Errorf("CreateUser() user = %v, want %v", user, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_DeleteAPIKey(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []Test[args, any]{
		{
			name: "delete Joe Biron's APIKey",
			args: args{
				ctx: context.Background(),
				id:  "2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := databaseRepository.DeleteAPIKey(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_DeleteUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []Test[args, bool]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := databaseRepository.DeleteUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetAPIKey(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []Test[args, *model.APIKey]{
		{
			name: "get Joe Bob's APIKey",
			args: args{
				ctx:    context.Background(),
				userId: "1",
			},
			want: &model.APIKey{
				Key: "12345abcdef",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotApiKey, err := databaseRepository.GetAPIKey(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotApiKey.Key, tt.want.Key) {
				t.Errorf("GetAPIKey() gotApiKey = %v, want %v", gotApiKey, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetById(t *testing.T) {
	type args struct {
		id int
	}
	type want struct {
		pronouns model.Pronouns
		exists   bool
	}
	tests := []Test[args, want]{
		{
			name: "get he/him",
			args: args{id: 1},
			want: want{
				pronouns: model.Pronouns{
					Subjective: "he",
					Objective:  "him",
				},
				exists: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pronouns, got1 := databaseRepository.GetById(tt.args.id)
			if !reflect.DeepEqual(pronouns, tt.want.pronouns) {
				t.Errorf("GetById() pronouns = %v, want %v", pronouns, tt.want)
			}
			if got1 != tt.want.exists {
				t.Errorf("GetById() got1 = %v, want %v", got1, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetByPronouns(t *testing.T) {
	type args struct {
		pronouns model.Pronouns
	}
	type want struct {
		id     int
		exists bool
	}
	tests := []Test[args, want]{
		{
			name: "get he him by pronouns",
			args: args{
				pronouns: model.Pronouns{
					Subjective: "he",
					Objective:  "him",
				},
			},
			want: want{
				id:     1,
				exists: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, exists := databaseRepository.GetByPronouns(tt.args.pronouns)
			if id != tt.want.id {
				t.Errorf("GetByPronouns() userId = %v, want %v", id, tt.want)
			}
			if exists != tt.want.exists {
				t.Errorf("GetByPronouns() exists = %v, want %v", exists, tt.want.exists)
			}
		})
	}
}

func TestDatabaseRepository_GetOAuth(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []Test[args, *model.OAuth]{
		{
			name: "get joe bob's oauth",
			args: args{
				ctx:    context.Background(),
				userId: "1",
			},
			want: &model.OAuth{
				Provider: models.ProviderGithub,
				UID:      "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := databaseRepository.GetOAuth(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOAuth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetOrCreatePronoun(t *testing.T) {
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		pronouns  model.Pronouns
		input     *model.NewUser
	}
	tests := []Test[args, *int]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := databaseRepository.GetOrCreatePronoun(tt.args.ctx, tt.args.queryable, tt.args.pronouns, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrCreatePronoun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrCreatePronoun() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetUserByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []Test[args, *model.User]{
		{
			name: "get Joe Bob by id",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			want: &model.User{
				ID:          "1",
				FirstName:   "Joe",
				LastName:    "Bob",
				Email:       "joe.bob@example.com",
				PhoneNumber: "100-200-3000",
				Pronouns: &model.Pronouns{
					Subjective: "he",
					Objective:  "him",
				},
				Age:               utils.Ptr(22),
				Role:              models.RoleNormal,
				Gender:            utils.Ptr("MALE"),
				Race:              []model.Race{"CAUCASIAN"},
				OAuth:             nil,
				MailingAddress:    nil,
				Mlh:               nil,
				ShirtSize:         model.ShirtSizeL,
				YearsOfExperience: utils.Ptr(3.5),
				EducationInfo:     nil,
				APIKey:            nil,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := databaseRepository.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetUserByOAuthUID(t *testing.T) {
	type args struct {
		ctx      context.Context
		oAuthUID string
		provider models.Provider
	}
	tests := []Test[args, *model.User]{
		{
			name: "get Joe Bob by oauth id",
			args: args{
				ctx:      context.Background(),
				oAuthUID: "1",
				provider: models.ProviderGithub,
			},
			want: &model.User{
				ID:          "1",
				FirstName:   "Joe",
				LastName:    "Bob",
				Email:       "joe.bob@example.com",
				PhoneNumber: "100-200-3000",
				Pronouns: &model.Pronouns{
					Subjective: "he",
					Objective:  "him",
				},
				Age:               utils.Ptr(22),
				Role:              models.RoleNormal,
				Gender:            utils.Ptr("MALE"),
				Race:              []model.Race{"CAUCASIAN"},
				OAuth:             nil,
				MailingAddress:    nil,
				Mlh:               nil,
				ShirtSize:         model.ShirtSizeL,
				YearsOfExperience: utils.Ptr(3.5),
				EducationInfo:     nil,
				APIKey:            nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := databaseRepository.GetUserByOAuthUID(tt.args.ctx, tt.args.oAuthUID, tt.args.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByOAuthUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByOAuthUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetUserMLHTerms(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []Test[args, *model.MLHTerms]{
		{
			name: "",
			args: args{
				ctx:    context.Background(),
				userId: "1",
			},
			want: &model.MLHTerms{
				SendMessages:  true,
				CodeOfConduct: true,
				ShareInfo:     true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := databaseRepository.GetUserMLHTerms(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserMLHTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserMLHTerms() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetUserMailingAddress(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []Test[args, *model.MailingAddress]{
		{
			name: "",
			args: args{
				ctx:    context.Background(),
				userId: "1",
			},
			want: &model.MailingAddress{
				Country:    "United States",
				State:      "Florida",
				City:       "Orlando",
				PostalCode: "32765",
				AddressLines: []string{
					"1000 Abc Rd",
					"APT 69",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := databaseRepository.GetUserMailingAddress(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserMailingAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserMailingAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_GetUsers(t *testing.T) {
	type args struct {
		ctx   context.Context
		first int
		after string
	}
	type want struct {
		users []*model.User
		total int
	}
	tests := []Test[args, want]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := databaseRepository.GetUsers(tt.args.ctx, tt.args.first, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(users, tt.want) {
				t.Errorf("GetUsers() users = %v, want %v", users, tt.want)
			}
			if total != tt.want.total {
				t.Errorf("GetUsers() total = %v, want %v", total, tt.want.total)
			}
		})
	}
}

func TestDatabaseRepository_InsertEducationInfo(t *testing.T) {
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		userId    int
		input     *model.EducationInfoInput
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := databaseRepository.InsertEducationInfo(tt.args.ctx, tt.args.queryable, tt.args.userId, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("InsertEducationInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_InsertMLHTerms(t *testing.T) {
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		userId    int
		input     *model.MLHTermsInput
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.InsertMLHTerms(tt.args.ctx, tt.args.queryable, tt.args.userId, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("InsertMLHTerms() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_InsertMailingAddress(t *testing.T) {
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		userId    int
		input     *model.MailingAddressInput
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := databaseRepository.InsertMailingAddress(tt.args.ctx, tt.args.queryable, tt.args.userId, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("InsertMailingAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_InsertUser(t *testing.T) {
	type args struct {
		ctx          context.Context
		queryable    shared_db_utils.Queryable
		input        *model.NewUser
		pronounIdPtr *int
		oAuth        *model.OAuth
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := databaseRepository.InsertUser(tt.args.ctx, tt.args.queryable, tt.args.input, tt.args.pronounIdPtr, tt.args.oAuth)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InsertUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_SearchUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []Test[args, []*model.User]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := databaseRepository.SearchUser(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(users, tt.want) {
				t.Errorf("SearchUser() users = %v, want %v", users, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_Set(t *testing.T) {
	type args struct {
		id       int
		pronouns model.Pronouns
	}
	tests := []struct {
		name string
		args args
	}{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			databaseRepository.Set(tt.args.id, tt.args.pronouns)
		})
	}
}

func TestDatabaseRepository_UpdateAge(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
		age *int
		tx  pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := databaseRepository.UpdateAge(tt.args.ctx, tt.args.id, tt.args.age, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateEducationInfo(t *testing.T) {

	type args struct {
		ctx   context.Context
		id    string
		input *model.EducationInfoUpdate
		tx    pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.UpdateEducationInfo(tt.args.ctx, tt.args.id, tt.args.input, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEducationInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateEmail(t *testing.T) {

	type args struct {
		ctx   context.Context
		id    string
		email *string
		tx    pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.UpdateEmail(tt.args.ctx, tt.args.id, tt.args.email, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateFirstName(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    string
		first *string
		tx    pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.UpdateFirstName(tt.args.ctx, tt.args.id, tt.args.first, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFirstName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateGender(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     string
		gender *string
		tx     pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := databaseRepository.UpdateGender(tt.args.ctx, tt.args.id, tt.args.gender, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGender() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateLastName(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   string
		last *string
		tx   pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.UpdateLastName(tt.args.ctx, tt.args.id, tt.args.last, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateLastName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateMLHTerms(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    string
		input *model.MLHTermsUpdate
		tx    pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.UpdateMLHTerms(tt.args.ctx, tt.args.id, tt.args.input, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMLHTerms() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateMailingAddress(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    string
		input *model.MailingAddressUpdate
		tx    pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := databaseRepository.UpdateMailingAddress(tt.args.ctx, tt.args.id, tt.args.input, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMailingAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdatePhoneNumber(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     string
		number *string
		tx     pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.UpdatePhoneNumber(tt.args.ctx, tt.args.id, tt.args.number, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdatePronouns(t *testing.T) {
	type args struct {
		ctx     context.Context
		id      string
		pronoun *model.PronounsInput
		tx      pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := databaseRepository.UpdatePronouns(tt.args.ctx, tt.args.id, tt.args.pronoun, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePronouns() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateRace(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    string
		races []model.Race
		tx    pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.UpdateRace(tt.args.ctx, tt.args.id, tt.args.races, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateShirtSize(t *testing.T) {
	type args struct {
		ctx       context.Context
		id        string
		shirtSize *model.ShirtSize
		tx        pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := databaseRepository.UpdateShirtSize(tt.args.ctx, tt.args.id, tt.args.shirtSize, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateShirtSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    string
		input *model.UpdatedUser
	}
	tests := []Test[args, *model.User]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := databaseRepository.UpdateUser(tt.args.ctx, tt.args.id, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_UpdateYearsOfExperience(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    string
		years *float64
		tx    pgx.Tx
	}
	tests := []Test[args, any]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := databaseRepository.UpdateYearsOfExperience(tt.args.ctx, tt.args.id, tt.args.years, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateYearsOfExperience() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_getUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []Test[args, *model.User]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := databaseRepository.GetUser(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_getUserWithTx(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		tx    pgx.Tx
		args  []interface{}
	}
	tests := []Test[args, *model.User]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := databaseRepository.GetUserWithTx(tt.args.ctx, tt.args.query, tt.args.tx, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserWithTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUserWithTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDatabaseRepository(t *testing.T) {
	type args struct {
		databasePool *pgxpool.Pool
	}
	tests := []Test[args, *database.DatabaseRepository]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := database.NewDatabaseRepository(context.Background(), tt.args.databasePool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDatabaseRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanUser(t *testing.T) {
	type args struct {
		user      *model.User
		scannable database.Scannable
	}
	tests := []Test[args, *int]{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := database.ScanUser(tt.args.user, tt.args.scannable)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
