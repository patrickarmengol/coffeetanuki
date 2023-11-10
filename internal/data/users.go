package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/patrickarmengol/coffeetanuki/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type password struct {
	plaintext *string // distinguish between "" and nil
	hash      []byte
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"create_at"`
	Version   int       `json:"-"`
}

type UserRepository struct {
	DB *sql.DB
}

// helpers

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// validate

func (u *User) Validate(v *validator.Validator) {
	// name
	v.CheckField(validator.NotBlank(u.Name), "name", "must be provided")
	v.CheckField(validator.MaxBytes(u.Name, 500), "name", "must not be more than 500 bytes long")

	// email
	v.CheckField(validator.NotBlank(u.Email), "email", "must be provided")
	v.CheckField(validator.Matches(u.Email, validator.EmailRX), "email", "must be a valid email address")

	// password
	if u.Password.plaintext != nil {
		v.CheckField(validator.NotBlank(*u.Password.plaintext), "password", "must be provided")
		v.CheckField(validator.MinBytes(*u.Password.plaintext, 8), "password", "must be at least 8 bytes long")
		v.CheckField(validator.MaxBytes(*u.Password.plaintext, 72), "password", "must be no more than 72 bytes long")
	}

	if u.Password.hash == nil {
		panic("missing password hash for user") // programmer error
	}
}

func (rep *UserRepository) Insert(user *User) error {
	stmt := `
	INSERT INTO users (name, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rep.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (rep *UserRepository) GetByEmail(email string) (*User, error) {
	stmt := `
	SELECT id, name, email, password_hash, activated, created_at, version
	FROM users
	WHERE email = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rep.DB.QueryRowContext(ctx, stmt, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.CreatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (rep *UserRepository) Update(user *User) error {
	stmt := `
	UPDATE users
	SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version
	`

	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rep.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
