package dba

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/patrickarmengol/coffeetanuki/internal/model"
)

// crud

// create

func CreateRoaster(ctx context.Context, dbtx DBTX, p *model.RoasterCreateParams) (*model.RoasterDB, error) {
	stmt := `
	INSERT INTO roasters (name, description, website, location)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{p.Name, p.Description, p.Website, p.Location}

	roaster := model.RoasterDB{
		Name:        p.Name,
		Description: p.Description,
		Website:     p.Website,
		Location:    p.Location,
	}

	err := dbtx.QueryRowContext(ctx, stmt, args...).Scan(&roaster.ID, &roaster.CreatedAt, &roaster.Version)
	if err != nil {
		return nil, err
	}

	return &roaster, nil
}

// read

func GetRoaster(ctx context.Context, dbtx DBTX, id int64) (*model.RoasterDB, error) {
	stmt := `
	SELECT id, name, description, website, location, created_at, version
	FROM roasters
	WHERE id = $1
	`

	args := []any{id}

	var roaster model.RoasterDB

	err := dbtx.QueryRowContext(ctx, stmt, args...).Scan(&roaster.ID, &roaster.Name, &roaster.Description, &roaster.Website, &roaster.Location, &roaster.CreatedAt, &roaster.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errRecordNotFound("roasters", id)
		default:
			return nil, err
		}
	}

	return &roaster, nil
}

func FindRoasters(ctx context.Context, dbtx DBTX, p *model.RoasterFilterParams) ([]*model.RoasterDB, error) {
	conditions := []string{}

	// search term will match if all space-delimited words are in the concatenation of searchable columns
	searchFields := []string{"name"}
	termConditions := fmt.Sprintf(`(CONCAT(%s) ILIKE ALL($1) OR $1 = '{}')`, strings.Join(searchFields, ", ' ', "))
	conditions = append(conditions, termConditions)

	wrappedWords := []string{}
	for _, w := range strings.Fields(p.SearchTerm) {
		wrappedWords = append(wrappedWords, fmt.Sprintf("%%%s%%", w))
	}
	wordArray := pq.Array(wrappedWords)

	stmt := fmt.Sprintf(`
		SELECT id, name, description, website, location, created_at, version
		FROM roasters
		WHERE %s
		ORDER BY %s %s, id ASC
	`, strings.Join(conditions, " AND "), p.SortField, p.SortDir)

	args := []any{wordArray}

	rows, err := dbtx.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roasters := []*model.RoasterDB{}
	for rows.Next() {
		var roaster model.RoasterDB

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

func UpdateRoaster(ctx context.Context, dbtx DBTX, p *model.RoasterEditParams) (*model.RoasterDB, error) {
	current, err := GetRoaster(ctx, dbtx, p.ID)
	if err != nil {
		return nil, err
	}

	stmt := `
	UPDATE roasters
	SET name = $3, description = $4, website = $5, location = $6, version = version + 1
	WHERE id = $1 AND version = $2
	RETURNING version
	`

	args := []any{current.ID, current.Version, p.Name, p.Description, p.Website, p.Location}

	roaster := model.RoasterDB{
		ID:          current.ID,
		Name:        p.Name,
		Description: p.Description,
		Website:     p.Website,
		Location:    p.Location,
		CreatedAt:   current.CreatedAt,
	}

	err = dbtx.QueryRowContext(ctx, stmt, args...).Scan(&roaster.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errEditConflict("roasters", roaster.ID)
		default:
			return nil, err
		}
	}

	return &roaster, nil
}

// delete

func DeleteRoaster(ctx context.Context, dbtx DBTX, id int64) error {
	stmt := `
	DELETE FROM roasters
	WHERE id = $1
	`

	result, err := dbtx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errRecordNotFound("roasters", id)
	}

	return nil
}

// association helpers

func AttachRoasterAssociations(ctx context.Context, dbtx DBTX, roaster *model.RoasterDB) error {
	beans, err := GetBeansForRoaster(ctx, dbtx, roaster.ID)
	if err != nil {
		return fmt.Errorf("attach roaster beans: %w", err)
	}
	roaster.Beans = beans

	return nil
}

func AttachManyRoasterAssociations(ctx context.Context, dbtx DBTX, roasters []*model.RoasterDB) error {
	// TODO: this seems like an exceedingly stupid way of doing this; should just left join
	beans, err := FindBeans(ctx, dbtx, &model.BeanFilterParams{SortField: "id", SortDir: "ASC"})
	if err != nil {
		return fmt.Errorf("attach roasters beans: %w", err)
	}

	rm := map[int64]*model.RoasterDB{}
	for _, r := range roasters {
		rm[r.ID] = r
	}

	for _, b := range beans {
		br, ok := rm[b.RoasterID]
		if !ok {
			return fmt.Errorf("attach roasters beans: %w", err)
		}
		br.Beans = append(br.Beans, b)
	}

	return nil
}
