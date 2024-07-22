package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"

	Model "github.com/Adhiana46/aegis-asesmen/internal/organization/data/model"
	DataSource "github.com/Adhiana46/aegis-asesmen/internal/organization/data/source"
	PkgDataSource "github.com/Adhiana46/aegis-asesmen/pkg/data_sources"
	"github.com/pkg/errors"
)

type organizationPersistentPostgres struct {
	db *PkgDataSource.PostgresDB
}

func NewOrganizationPersistentPostgres(db *PkgDataSource.PostgresDB) DataSource.IOrganizationPersistent {
	return &organizationPersistentPostgres{
		db: db,
	}
}

func (u *organizationPersistentPostgres) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "organizationPersistentPostgres:" + chunks[len(chunks)-1]
}

func (r *organizationPersistentPostgres) GetList(ctx context.Context, offset, limit int) ([]*Model.Organization, error) {
	path := r.path()

	query := `
		SELECT
			id,
			name,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM "organizations"
	`

	aWheres := []string{}
	aOrders := []string{}
	args := []any{}

	// Filters
	if len(aWheres) > 0 {
		query += " WHERE " + strings.Join(aWheres, " AND ")
	}

	// Orders
	if len(aOrders) > 0 {
		query += " ORDER BY " + strings.Join(aOrders, ", ")
	}

	// offset
	query += fmt.Sprintf(" OFFSET %v", offset)

	// limit
	if limit != -1 {
		query += fmt.Sprintf(" Limit %v", limit)
	}

	models := []*Model.Organization{}
	err := r.db.SelectContext(ctx, &models, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}

	return models, nil
}

func (r *organizationPersistentPostgres) CountList(ctx context.Context) (int, error) {
	path := r.path()

	query := `
		SELECT
			COUNT(id) AS numrows
		FROM "organizations"
	`

	aWheres := []string{}
	args := []any{}

	// Filters
	if len(aWheres) > 0 {
		query += " WHERE " + strings.Join(aWheres, " AND ")
	}

	numrows := 0
	result, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, path)
	}

	if result.Next() {
		if err := result.Scan(&numrows); err != nil {
			return 0, errors.Wrap(err, path)
		}
	}

	return numrows, nil
}

func (r *organizationPersistentPostgres) GetByID(ctx context.Context, id string) (*Model.Organization, error) {
	path := r.path()

	query := `
		SELECT
			id,
			name,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM "organizations"
		WHERE id = $1
	`

	model := Model.Organization{}
	err := r.db.GetContext(ctx, &model, query, id)
	if err == nil {
		return &model, nil
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return nil, errors.Wrap(err, path)
}

func (r *organizationPersistentPostgres) GetByName(ctx context.Context, name string) (*Model.Organization, error) {
	path := r.path()

	query := `
		SELECT
			id,
			name,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM "organizations"
		WHERE name = $1
	`

	model := Model.Organization{}
	err := r.db.GetContext(ctx, &model, query, name)
	if err == nil {
		return &model, nil
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return nil, errors.Wrap(err, path)
}

func (r *organizationPersistentPostgres) Store(ctx context.Context, model *Model.Organization) error {
	path := r.path()

	query := `
		INSERT INTO "organizations"
		(id, name, created_at, created_by, updated_at, updated_by)
		VALUES
		(:id, :name, :created_at, :created_by, :updated_at, :updated_by)
	`

	_, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}

func (r *organizationPersistentPostgres) Update(ctx context.Context, model *Model.Organization) error {
	path := r.path()

	query := `
        UPDATE "organizations"
        SET
            name = :name,
			created_at = :created_at,
			created_by = :created_by,
			updated_at = :updated_at,
			updated_by = :updated_by
        WHERE
            id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}

func (r *organizationPersistentPostgres) Destroy(ctx context.Context, model *Model.Organization) error {
	path := r.path()

	query := `
		DELETE FROM "organizations" WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, model.Id)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}
