package dba

import (
	"context"

	"github.com/lib/pq"
	"github.com/patrickarmengol/somethingsomethingcoffee/internal/model"
)

// TODO: maybe write a setRolePermisisonsForUser() to set based on roles; define roles in model package

func AddPermissionsForUser(ctx context.Context, dbtx DBTX, id int64, codes model.PermissionCodes) error {
	stmt := `
	INSERT INTO users_permissions
	SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
	`

	_, err := dbtx.ExecContext(ctx, stmt, id, pq.Array(codes))
	return err
}

func GetPermissionsForUser(ctx context.Context, dbtx DBTX, id int64) (model.PermissionCodes, error) {
	stmt := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
	INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = $1
	`

	rows, err := dbtx.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pcs model.PermissionCodes

	for rows.Next() {
		var pc string
		err := rows.Scan(&pc)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, pc)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pcs, nil
}
