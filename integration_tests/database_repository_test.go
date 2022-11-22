package integration_tests

import (
	"context"
	"flag"
	"fmt"
	shared_db_utils "github.com/KnightHacks/knighthacks_shared/database"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"reflect"
	"testing"
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

	databaseRepository = database.NewDatabaseRepository(pool)
	os.Exit(t.Run())
}

func TestDatabaseRepository_AddAPIKey(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.APIKey
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.AddAPIKey(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddAPIKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_CreateUser(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		oAuth *model.OAuth
		input *model.NewUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.CreateUser(tt.args.ctx, tt.args.oAuth, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_DeleteAPIKey(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.DeleteAPIKey(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_DeleteUser(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.DeleteUser(tt.args.ctx, tt.args.id)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx context.Context
		obj *model.User
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantApiKey *model.APIKey
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			gotApiKey, err := r.GetAPIKey(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotApiKey, tt.wantApiKey) {
				t.Errorf("GetAPIKey() gotApiKey = %v, want %v", gotApiKey, tt.wantApiKey)
			}
		})
	}
}

func TestDatabaseRepository_GetById(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		id int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.Pronouns
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, got1 := r.GetById(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetById() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetById() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDatabaseRepository_GetByPronouns(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		pronouns model.Pronouns
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, got1 := r.GetByPronouns(tt.args.pronouns)
			if got != tt.want {
				t.Errorf("GetByPronouns() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetByPronouns() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDatabaseRepository_GetOAuth(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.OAuth
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.GetOAuth(tt.args.ctx, tt.args.id)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		pronouns  model.Pronouns
		input     *model.NewUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.GetOrCreatePronoun(tt.args.ctx, tt.args.queryable, tt.args.pronouns, tt.args.input)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.GetUserByID(tt.args.ctx, tt.args.id)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx      context.Context
		oAuthUID string
		provider models.Provider
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.GetUserByOAuthUID(tt.args.ctx, tt.args.oAuthUID, tt.args.provider)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.MLHTerms
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.GetUserMLHTerms(tt.args.ctx, tt.args.userId)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.MailingAddress
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.GetUserMailingAddress(tt.args.ctx, tt.args.userId)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		first int
		after string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.User
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, got1, err := r.GetUsers(tt.args.ctx, tt.args.first, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUsers() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetUsers() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDatabaseRepository_InsertEducationInfo(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		userId    int
		input     *model.EducationInfoInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.InsertEducationInfo(tt.args.ctx, tt.args.queryable, tt.args.userId, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("InsertEducationInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_InsertMLHTerms(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		userId    int
		input     *model.MLHTermsInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.InsertMLHTerms(tt.args.ctx, tt.args.queryable, tt.args.userId, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("InsertMLHTerms() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_InsertMailingAddress(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		userId    int
		input     *model.MailingAddressInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.InsertMailingAddress(tt.args.ctx, tt.args.queryable, tt.args.userId, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("InsertMailingAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_InsertUser(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx          context.Context
		queryable    shared_db_utils.Queryable
		input        *model.NewUser
		pronounIdPtr *int
		oAuth        *model.OAuth
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.InsertUser(tt.args.ctx, tt.args.queryable, tt.args.input, tt.args.pronounIdPtr, tt.args.oAuth)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.SearchUser(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseRepository_Set(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		id       int
		pronouns model.Pronouns
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			r.Set(tt.args.id, tt.args.pronouns)
		})
	}
}

func TestDatabaseRepository_UpdateAge(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx context.Context
		id  string
		age *int
		tx  pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateAge(tt.args.ctx, tt.args.id, tt.args.age, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateEducationInfo(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		id    string
		input *model.EducationInfoUpdate
		tx    pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateEducationInfo(tt.args.ctx, tt.args.id, tt.args.input, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEducationInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateEmail(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		id    string
		email *string
		tx    pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateEmail(tt.args.ctx, tt.args.id, tt.args.email, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateFirstName(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		id    string
		first *string
		tx    pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateFirstName(tt.args.ctx, tt.args.id, tt.args.first, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFirstName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateGender(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx    context.Context
		id     string
		gender *string
		tx     pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateGender(tt.args.ctx, tt.args.id, tt.args.gender, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGender() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateLastName(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx  context.Context
		id   string
		last *string
		tx   pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateLastName(tt.args.ctx, tt.args.id, tt.args.last, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateLastName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateMLHTerms(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		id    string
		input *model.MLHTermsUpdate
		tx    pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateMLHTerms(tt.args.ctx, tt.args.id, tt.args.input, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMLHTerms() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateMailingAddress(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		id    string
		input *model.MailingAddressUpdate
		tx    pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateMailingAddress(tt.args.ctx, tt.args.id, tt.args.input, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMailingAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdatePhoneNumber(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx    context.Context
		id     string
		number *string
		tx     pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdatePhoneNumber(tt.args.ctx, tt.args.id, tt.args.number, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdatePronouns(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx     context.Context
		id      string
		pronoun *model.PronounsInput
		tx      pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdatePronouns(tt.args.ctx, tt.args.id, tt.args.pronoun, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePronouns() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateRace(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		id    string
		races []model.Race
		tx    pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateRace(tt.args.ctx, tt.args.id, tt.args.races, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateShirtSize(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx       context.Context
		id        string
		shirtSize *model.ShirtSize
		tx        pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateShirtSize(tt.args.ctx, tt.args.id, tt.args.shirtSize, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateShirtSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_UpdateUser(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		id    string
		input *model.UpdatedUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.UpdateUser(tt.args.ctx, tt.args.id, tt.args.input)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		id    string
		years *float64
		tx    pgx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.UpdateYearsOfExperience(tt.args.ctx, tt.args.id, tt.args.years, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateYearsOfExperience() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_getPronouns(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx       context.Context
		queryable shared_db_utils.Queryable
		pronounId int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			if err := r.GetPronouns(tt.args.ctx, tt.args.queryable, tt.args.pronounId); (err != nil) != tt.wantErr {
				t.Errorf("getPronouns() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseRepository_getUser(t *testing.T) {
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.GetUser(tt.args.ctx, tt.args.query, tt.args.args...)
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
	type fields struct {
		DatabasePool      *pgxpool.Pool
		PronounMap        map[int]model.Pronouns
		PronounReverseMap map[model.Pronouns]int
	}
	type args struct {
		ctx   context.Context
		query string
		tx    pgx.Tx
		args  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := database.DatabaseRepository{
				DatabasePool:      tt.fields.DatabasePool,
				PronounMap:        tt.fields.PronounMap,
				PronounReverseMap: tt.fields.PronounReverseMap,
			}
			got, err := r.GetUserWithTx(tt.args.ctx, tt.args.query, tt.args.tx, tt.args.args...)
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

func TestGenerateAPIKey(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := database.GenerateAPIKey(tt.args.length); got != tt.want {
				t.Errorf("GenerateAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDatabaseRepository(t *testing.T) {
	type args struct {
		databasePool *pgxpool.Pool
	}
	tests := []struct {
		name string
		args args
		want *database.DatabaseRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := database.NewDatabaseRepository(tt.args.databasePool); !reflect.DeepEqual(got, tt.want) {
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
	tests := []struct {
		name    string
		args    args
		want    *int
		wantErr bool
	}{
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
