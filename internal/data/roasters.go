package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Roaster struct {
	ID          int64
	Name        string
	Description string
	Website     string
	Location    string
	CreatedAt   time.Time
	Version     int
}

type RoasterModel struct {
	DB *sql.DB
}

// create

func (m *RoasterModel) Insert(roaster *Roaster) error {
	stmt := `
	INSERT INTO roasters (name, description, website, location)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{roaster.Name, roaster.Description, roaster.Website, roaster.Location}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, stmt, args...).Scan(&roaster.ID, &roaster.CreatedAt, &roaster.Version)
}

// read

func (m *RoasterModel) Get(id int64) (*Roaster, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	stmt := `
	SELECT id, name, description, website, location, created_at, version
	FROM roasters
	WHERE id = $1
	`

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var roaster Roaster
	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&roaster.ID, &roaster.Name, &roaster.Description, &roaster.Website, &roaster.Location, &roaster.CreatedAt, &roaster.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &roaster, nil
}

func (m *RoasterModel) GetAll() ([]*Roaster, error) {
	stmt := `
	SELECT id, name, description, website, location, created_at, version
	FROM roasters
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roasters := []*Roaster{}
	for rows.Next() {
		var roaster Roaster

		err := rows.Scan(&roaster.ID, &roaster.Name, &roaster.Description, &roaster.Website, &roaster.Location, &roaster.CreatedAt, &roaster.Version)
		if err != nil {
			return nil, err
		}

		roasters = append(roasters, &roaster)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roasters, nil
}

// update

func (m *RoasterModel) Update(roaster *Roaster) error {
	stmt := `
	UPDATE roasters
	SET name = $1, description = $1, website = $1, location = $2, version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING version
	`

	args := []any{roaster.Name, roaster.Description, roaster.Website, roaster.Location, roaster.ID, roaster.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&roaster.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// delete

func (m *RoasterModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	stmt := `
	DELETE FROM roasters
	WHERE id = $1
	`

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
