package model

import (
	"time"

	"github.com/patrickarmengol/coffeetanuki/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &UserResponse{}

type UserCreateInput struct {
	Name              string `form:"name"`
	Email             string `form:"email"`
	PasswordPlaintext string `form:"password"`

	validator.Validator `form:"-"`
}

func (i *UserCreateInput) Validate() {
	i.CheckField(validator.NotBlank(i.Name), "name", "this field cannot be blank")
	i.CheckField(validator.MaxChars(i.Name, 20), "name", "this field must be at most 20 characters")
	i.CheckField(validator.NotBlank(i.Email), "email", "this field cannot be blank")
	i.CheckField(validator.Matches(i.Email, validator.EmailRX), "email", "this field must be a valid email")
	i.CheckField(validator.NotBlank(i.PasswordPlaintext), "password", "this field cannot be blank")
	i.CheckField(validator.MinChars(i.PasswordPlaintext, 8), "password", "this field must be at least 8 characters")
	i.CheckField(validator.MaxChars(i.PasswordPlaintext, 30), "password", "this field must be at most 30 characters")
	i.CheckField(validator.MaxBytes(i.PasswordPlaintext, 72), "password", "this field must be at most 72 bytes")
}

func (i *UserCreateInput) ToParams() (*UserCreateParams, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(i.PasswordPlaintext), 12)
	if err != nil {
		return nil, err
	}

	return &UserCreateParams{
		Name:         i.Name,
		Email:        i.Email,
		PasswordHash: hash,
		Activated:    true, // TODO: change to false after email registration confirmation implemented
	}, nil
}

type UserCreateParams struct {
	Name         string
	Email        string
	PasswordHash []byte
	Activated    bool
}

type UserLoginInput struct {
	Email             string `form:"email"`
	PasswordPlaintext string `form:"password"`

	validator.Validator `form:"-"`
}

func (i *UserLoginInput) Validate() {
	i.CheckField(validator.NotBlank(i.Email), "email", "this field cannot be blank")
	i.CheckField(validator.Matches(i.Email, validator.EmailRX), "email", "this field must be a valid email")
	i.CheckField(validator.NotBlank(i.PasswordPlaintext), "password", "this field cannot be blank")
	i.CheckField(validator.MinChars(i.PasswordPlaintext, 8), "password", "this field must be at least 8 characters")
	i.CheckField(validator.MaxChars(i.PasswordPlaintext, 30), "password", "this field must be at most 30 characters")
	i.CheckField(validator.MaxBytes(i.PasswordPlaintext, 72), "password", "this field must be at most 72 bytes")
}

type UserDB struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash []byte
	Activated    bool
	CreatedAt    time.Time
	Version      int
}

func (m *UserDB) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Activated: m.Activated,
	}
}

type UserResponse struct {
	ID        int64
	Name      string
	Email     string
	Activated bool
}

func (r *UserResponse) IsAnonymous() bool {
	return r == AnonymousUser
}
