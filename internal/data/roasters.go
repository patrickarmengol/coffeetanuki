package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/patrickarmengol/coffeetanuki/internal/validator"
)

type Roaster struct {
	ID          int64
	Name        string
	Description string
	Website     string
	Location    string
	Beans       []*Bean
	CreatedAt   time.Time
	Version     int
}

type RoasterRepository struct {
	DB *sql.DB
}

// validate

func (r *Roaster) Validate(v *validator.Validator) {
	v.CheckField(validator.NotBlank(r.Name), "name", "This field cannot be blank")
	v.CheckField(validator.NotBlank(r.Description), "description", "This field cannot be blank")
	v.CheckField(validator.IsURL(r.Website), "website", "This field must be a valid URL (eg. https://coffeetanuki.com)")
	v.CheckField(validator.Matches(r.Location, validator.LocationRX), "location", "This field must be a valid location (eg. Seattle, Washington, USA)")
}

// create

func (rep RoasterRepository) Insert(roaster *Roaster) error {
	stmt := `
	INSERT INTO roasters (name, description, website, location)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{roaster.Name, roaster.Description, roaster.Website, roaster.Location}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return rep.DB.QueryRowContext(ctx, stmt, args...).Scan(&roaster.ID, &roaster.CreatedAt, &roaster.Version)
}

// read

func (rep RoasterRepository) Get(id int64) (*Roaster, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	stmt := `
	SELECT *
	FROM roasters
	WHERE id = $1
	`

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var roaster Roaster
	err := rep.DB.QueryRowContext(ctx, stmt, args...).Scan(
		&roaster.ID,
		&roaster.Name,
		&roaster.Description,
		&roaster.Website,
		&roaster.Location,
		&roaster.CreatedAt,
		&roaster.Version,
	)
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

func (rep RoasterRepository) GetFull(id int64) (*Roaster, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	stmt := `
	SELECT r.*, b.*
	FROM roasters r
	LEFT JOIN beans b ON r.id = b.roaster_id
	WHERE r.id = $1
	`

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := rep.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roaster Roaster

	for rows.Next() {
		var nbean NullableBean

		err := rows.Scan(
			&roaster.ID,
			&roaster.Name,
			&roaster.Description,
			&roaster.Website,
			&roaster.Location,
			&roaster.CreatedAt,
			&roaster.Version,
			&nbean.ID,
			&nbean.Name,
			&nbean.RoastLevel,
			&nbean.RoasterID,
			&nbean.CreatedAt,
			&nbean.Version,
		)
		if err != nil {
			return nil, err
		}

		if nbean.ID.Valid {
			bean := nbean.UnNullify()
			roaster.Beans = append(roaster.Beans, &bean)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &roaster, nil
}

// func (rep RoasterRepository) GetAll() ([]*Roaster, error) {
// 	stmt := `
// 	SELECT *
// 	FROM roasters
// 	`
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
//
// 	rows, err := rep.DB.QueryContext(ctx, stmt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	roasters := []*Roaster{}
//
// 	for rows.Next() {
// 		var roaster Roaster
//
// 		err := rows.Scan(
// 			&roaster.ID,
// 			&roaster.Name,
// 			&roaster.Description,
// 			&roaster.Website,
// 			&roaster.Location,
// 			&roaster.CreatedAt,
// 			&roaster.Version,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		roasters = append(roasters, &roaster)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
//
// 	return roasters, nil
// }

func (rep RoasterRepository) GetAllFull() ([]*Roaster, error) {
	stmt := `
	SELECT r.*, b.*
	FROM roasters r
	LEFT JOIN beans b
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := rep.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roastersMap := make(map[int64]Roaster)

	for rows.Next() {
		var roaster Roaster
		var nbean NullableBean

		err := rows.Scan(
			&roaster.ID,
			&roaster.Name,
			&roaster.Description,
			&roaster.Website,
			&roaster.Location,
			&roaster.CreatedAt,
			&roaster.Version,
			&nbean.ID,
			&nbean.Name,
			&nbean.RoastLevel,
			&nbean.RoasterID,
			&nbean.CreatedAt,
			&nbean.Version,
		)
		if err != nil {
			return nil, err
		}

		rr, ok := roastersMap[roaster.ID]
		if !ok {
			rr = roaster
		}

		if nbean.ID.Valid {
			bean := nbean.UnNullify()
			rr.Beans = append(rr.Beans, &bean)
		}

		roastersMap[roaster.ID] = rr

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	roasters := make([]*Roaster, 0, len(roastersMap))
	for _, r := range roastersMap {
		roasters = append(roasters, &r)
	}
	return roasters, nil
}

func (rep RoasterRepository) Search(sq SearchQuery) ([]*Roaster, error) {
	// TODO: clean this up

	conditions := []string{}

	// search term will match if all space-delimited words are in the concatenation of searchable columns
	searchFields := []string{"name", "description"}
	termConditions := fmt.Sprintf(`CONCAT(%s) ILIKE ALL($1)`, strings.Join(searchFields, ", ' ', "))
	conditions = append(conditions, termConditions)

	// TODO: move this to SearchQuery method
	wrappedWords := []string{}
	for _, w := range strings.Fields(sq.Term) {
		wrappedWords = append(wrappedWords, fmt.Sprintf("%%%s%%", w))
	}
	wordArray := pq.Array(wrappedWords)

	// TODO: add tag conditions

	stmt := fmt.Sprintf(`
		SELECT id, name, description, website, location, created_at, version
		FROM roasters
		WHERE (%s) OR $1 = '{}'
		ORDER BY %s %s, id ASC
	`, strings.Join(conditions, " AND "), sq.sortBy(), sq.sortDir())

	args := []any{wordArray}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := rep.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roasters := []*Roaster{}

	for rows.Next() {
		var roaster Roaster

		err := rows.Scan(
			&roaster.ID,
			&roaster.Name,
			&roaster.Description,
			&roaster.Website,
			&roaster.Location,
			&roaster.CreatedAt,
			&roaster.Version,
		)
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

func (rep RoasterRepository) Update(roaster *Roaster) error {
	stmt := `
	UPDATE roasters
	SET name = $1, description = $2, website = $3, location = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version
	`

	args := []any{roaster.Name, roaster.Description, roaster.Website, roaster.Location, roaster.ID, roaster.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rep.DB.QueryRowContext(ctx, stmt, args...).Scan(&roaster.Version)
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

func (rep RoasterRepository) Delete(id int64) error {
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

	result, err := rep.DB.ExecContext(ctx, stmt, args...)
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
