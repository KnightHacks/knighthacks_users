// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type LoginPayload struct {
	// If false then you must register immediately following this. Else, you are logged in and have access to your own user.
	AccountExists bool    `json:"accountExists"`
	User          *User   `json:"user"`
	Jwt           *string `json:"jwt"`
}

type NewUser struct {
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phoneNumber"`
	Pronouns    *PronounsInput `json:"pronouns"`
	Age         *int           `json:"age"`
}

type OAuth struct {
	Provider    Provider `json:"provider"`
	AccessToken string   `json:"accessToken"`
}

// Example:
// subjective=he
// objective=him
type Pronouns struct {
	Subjective string `json:"subjective"`
	Objective  string `json:"objective"`
}

type PronounsInput struct {
	SubjectivePersonal string `json:"subjectivePersonal"`
	ObjectivePersonal  string `json:"objectivePersonal"`
	Reflexive          string `json:"reflexive"`
}

type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	FullName    string    `json:"fullName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Pronouns    *Pronouns `json:"pronouns"`
	Age         *int      `json:"age"`
	OAuth       *OAuth    `json:"oAuth"`
}

func (User) IsEntity() {}

type Provider string

const (
	ProviderGithub Provider = "GITHUB"
	ProviderGmail  Provider = "GMAIL"
)

var AllProvider = []Provider{
	ProviderGithub,
	ProviderGmail,
}

func (e Provider) IsValid() bool {
	switch e {
	case ProviderGithub, ProviderGmail:
		return true
	}
	return false
}

func (e Provider) String() string {
	return string(e)
}

func (e *Provider) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Provider(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Provider", str)
	}
	return nil
}

func (e Provider) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
