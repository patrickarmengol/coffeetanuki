package dba

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/patrickarmengol/somethingsomethingcoffee/internal/model"
)

// crud

// create

func CreateBean(ctx context.Context, dbtx DBTX, p *model.BeanCreateParams) (*model.BeanDB, error) {
	stmt := `
	INSERT INTO beans (name, roast_level, roaster_id)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version
	`

	args := []any{p.Name, p.RoastLevel, p.RoasterID}

	bean := model.BeanDB{
		Name:       p.Name,
		RoastLevel: p.RoastLevel,
		RoasterID:  p.RoasterID,
	}

	err := dbtx.QueryRowContext(ctx, stmt, args...).Scan(&bean.ID, &bean.CreatedAt, &bean.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "beans" violates foreign key constraint "beans_roaster_id_fkey"`:
			return nil, errInvalidFK("beans", "roaster_id", p.RoasterID)
		default:
			return nil, err
		}
	}

	return &bean, nil
}

// read

func GetBean(ctx context.Context, dbtx DBTX, id int64) (*model.BeanDB, error) {
	stmt := `
	SELECT id, name, roast_level, roaster_id, created_at, version
	FROM beans
	WHERE id = $1
	`

	args := []any{id}

	var bean model.BeanDB

	err := dbtx.QueryRowContext(ctx, stmt, args...).Scan(&bean.ID, &bean.Name, &bean.RoastLevel, &bean.RoasterID, &bean.CreatedAt, &bean.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errRecordNotFound("beans", id)
		default:
			return nil, err
		}
	}

	return &bean, nil
}

func FindBeans(ctx context.Context, dbtx DBTX, p *model.BeanFilterParams) ([]*model.BeanDB, error) {
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
		SELECT id, name, roast_level, roaster_id, created_at, version
		FROM beans
		WHERE %s
		ORDER BY %s %s, id ASC
	`, strings.Join(conditions, " AND "), p.SortField, p.SortDir)

	args := []any{wordArray}

	rows, err := dbtx.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beans := []*model.BeanDB{}
	for rows.Next() {
		var bean model.BeanDB

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

func UpdateBean(ctx context.Context, dbtx DBTX, p *model.BeanEditParams) (*model.BeanDB, error) {
	current, err := GetBean(ctx, dbtx, p.ID)
	if err != nil {
		return nil, err
	}

	stmt := `
	UPDATE beans
	SET name = $3, roast_level = $4, roaster_id = $5, version = version + 1
	WHERE id = $1 AND version = $2
	RETURNING version
	`

	args := []any{current.ID, current.Version, p.Name, p.RoastLevel, p.RoasterID}

	bean := model.BeanDB{
		ID:         current.ID,
		Name:       p.Name,
		RoastLevel: p.RoastLevel,
		RoasterID:  p.RoasterID,
		CreatedAt:  current.CreatedAt,
	}

	err = dbtx.QueryRowContext(ctx, stmt, args...).Scan(&bean.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "beans" violates foreign key constraint "beans_roaster_id_fkey"`:
			return nil, errInvalidFK("beans", "roaster_id", p.RoasterID)
		case errors.Is(err, sql.ErrNoRows):
			return nil, errEditConflict("beans", bean.ID)
		default:
			return nil, err
		}
	}

	return &bean, nil
}

// delete

func DeleteBean(ctx context.Context, dbtx DBTX, id int64) error {
	stmt := `
	DELETE FROM beans
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
		return errRecordNotFound("beans", id)
	}

	return nil
}

// special

// TODO: move this functionality into FindBeans
func GetBeansForRoaster(ctx context.Context, dbtx DBTX, id int64) ([]*model.BeanDB, error) {
	stmt := `
	SELECT id, name, roast_level, roaster_id, created_at, version
	FROM beans
	WHERE roaster_id = $1
	`

	args := []any{id}

	rows, err := dbtx.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}

	beans := []*model.BeanDB{}
	for rows.Next() {
		var bean model.BeanDB

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

// association helpers

func AttachBeanAssociations(ctx context.Context, dbtx DBTX, bean *model.BeanDB) error {
	roaster, err := GetRoaster(ctx, dbtx, bean.RoasterID)
	if err != nil {
		return fmt.Errorf("attach bean roaster: %w", err)
	}
	bean.Roaster = roaster

	return nil
}

func AttachManyBeanAssociations(ctx context.Context, dbtx DBTX, beans []*model.BeanDB) error {
	// TODO: create a filter for IN() some set of roaster_ids; avoid retrieving full table
	// how do i get a unique set of roaster_ids?
	roasters, err := FindRoasters(ctx, dbtx, &model.RoasterFilterParams{SortField: "id", SortDir: "ASC"})
	if err != nil {
		return fmt.Errorf("attach beans roaster: %w", err)
	}

	rm := map[int64]*model.RoasterDB{}
	for _, r := range roasters {
		rm[r.ID] = r
	}

	for _, b := range beans {
		br, ok := rm[b.RoasterID]
		if !ok {
			return fmt.Errorf("attach beans roaster: %w", err)
		}
		b.Roaster = br
	}

	return nil
}
