package tableaccess

import (
	"context"
	"database/sql"
	"fmt"
)

type Manager interface {
}

type manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) Manager {
	return &manager{db: db}
}

func (m *manager) UserCanAccessTable(ctx context.Context, username, table string) (bool, error) {
	results, err := m.db.QueryContext(ctx, fmt.Sprintf(`SELECT 
    grantee,
    table_schema,
    table_name,
    privilege_type,
    grantor
FROM information_schema.table_privileges
WHERE grantee = '%s'`, username))
	if err != nil {
		return false, err
	}

	for results.Next() {
		var (
			grantee,
			tableSchema,
			tableName,
			privilegeType,
			grantor string
		)
		if err = results.Scan(
			&grantee,
			&tableSchema,
			&tableName,
			&privilegeType,
			&grantor,
		); err != nil {
			return false, err
		}

		fmt.Printf(`
grantee: %s
tableSchema: %s
tableName: %s
privilegeType: %s
grantor: %s
`, grantee, tableSchema, tableName, privilegeType, grantor)
	}

	return false, nil
}
