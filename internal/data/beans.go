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

var ErrInvalidRoasterID = errors.New("invalid roaster id")

var roastLevels = []string{"light", "medium-light", "medium", "medium-dark", "dark"}

type Bean struct {
	ID         int64
	Name       string
	RoastLevel string
	RoasterID  int64
	CreatedAt  time.Time
	Version    int
}

type NullableBean struct {
	ID         sql.NullInt64
	Name       sql.NullString
	RoastLevel sql.NullString
	RoasterID  sql.NullInt64
	CreatedAt  sql.NullTime
	Version    sql.NullInt32
}

type BeanRepository struct {
	DB *sql.DB
}

// validate

func (b *Bean) Validate(v *validator.Validator) {
	v.CheckField(validator.NotBlank(b.Name), "name", "This field cannot be blank")
	v.CheckField(validator.NotBlank(b.RoastLevel), "roast_level", "This field cannot be blank")
	v.CheckField(validator.PermittedValue(b.RoastLevel, roastLevels...), "roast_level", fmt.Sprintf("This field must be one of %v", roastLevels))
}

// helpers

func (nb NullableBean) UnNullify() Bean {
	return Bean{
		ID:         nb.ID.Int64,
		Name:       nb.Name.String,
		RoastLevel: nb.RoastLevel.String,
		RoasterID:  nb.RoasterID.Int64,
		CreatedAt:  nb.CreatedAt.Time,
		Version:    int(nb.Version.Int32),
	}
}

// create

func (rep BeanRepository) Insert(bean *Bean) error {
	stmt := `
	INSERT INTO beans (name, roast_level, roaster_id)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version
	`

	args := []any{bean.Name, bean.RoastLevel, bean.RoasterID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rep.DB.QueryRowContext(ctx, stmt, args...).Scan(&bean.ID, &bean.CreatedAt, &bean.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "beans" violates foreign key constraint "beans_roaster_id_fkey"`:
			return ErrInvalidRoasterID
		default:
			return err
		}
	}

	return nil
}

// read

func (rep BeanRepository) Get(id int64) (*Bean, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	stmt := `
	SELECT id, name, roast_level, roaster_id, created_at, version
	FROM beans
	WHERE id = $1
	`

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var bean Bean
	err := rep.DB.QueryRowContext(ctx, stmt, args...).Scan(&bean.ID, &bean.Name, &bean.RoastLevel, &bean.RoasterID, &bean.CreatedAt, &bean.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &bean, nil
}

func (rep BeanRepository) GetAll() ([]*Bean, error) {
	stmt := `
	SELECT id, name, roast_level, roaster_id, created_at, version
	FROM beans
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := rep.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beans := []*Bean{}
	for rows.Next() {
		var bean Bean

		err := rows.Scan(&bean.ID, &bean.Name, &bean.RoastLevel, &bean.RoasterID, &bean.CreatedAt, &bean.Version)
		if err != nil {
			return nil, err
		}

		beans = append(beans, &bean)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return beans, nil
}

func (rep BeanRepository) GetAllForRoaster(roasterID int64) ([]*Bean, error) {
	stmt := `
	SELECT id, name, roast_level, created_at, version
	FROM roasters
	WHERE roaster_id = $1
	`

	args := []any{roasterID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := rep.DB.QueryContext(ctx, stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beans := []*Bean{}
	for rows.Next() {
		var bean Bean

		err := rows.Scan(&bean.ID, &bean.Name, &bean.RoastLevel, &bean.CreatedAt, &bean.Version)
		if err != nil {
			return nil, err
		}

		beans = append(beans, &bean)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return beans, nil
}

func (rep BeanRepository) Search(sq SearchQuery) ([]*Bean, error) {
	conditions := []string{}

	// search term will match if all space-delimited words are in the concatenation of searchable columns
	searchFields := []string{"name"}
	termConditions := fmt.Sprintf(`(CONCAT(%s) ILIKE ALL($1) OR $1 = '{}')`, strings.Join(searchFields, ", ' ', "))
	conditions = append(conditions, termConditions)

	wordArray := pq.Array(sq.termWordsWrapped())

	stmt := fmt.Sprintf(`
		SELECT id, name, roast_level, roaster_id, created_at, version
		FROM beans
		WHERE %s
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

	beans := []*Bean{}
	for rows.Next() {
		var bean Bean

		err := rows.Scan(&bean.ID, &bean.Name, &bean.RoastLevel, &bean.RoasterID, &bean.CreatedAt, &bean.Version)
		if err != nil {
			return nil, err
		}

		beans = append(beans, &bean)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return beans, nil
}

// update

func (rep BeanRepository) Update(bean *Bean) error {
	stmt := `
	UPDATE beans
	SET name = $1, roast_level = $2, roaster_id = $3, version = version + 1
	WHERE id = $4 AND version = $5
	RETURNING version
	`

	args := []any{bean.Name, bean.RoastLevel, bean.RoasterID, bean.ID, bean.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rep.DB.QueryRowContext(ctx, stmt, args...).Scan(&bean.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "beans" violates foreign key constraint "beans_roaster_id_fkey"`:
			return ErrInvalidRoasterID
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// delete

func (rep BeanRepository) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM beans
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := rep.DB.ExecContext(ctx, query, id)
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
