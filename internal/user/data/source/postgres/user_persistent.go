package postgres

import (
	"context"
	"database/sql"

	Model "github.com/Adhiana46/aegis-asesmen/internal/user/data/model"
	Source "github.com/Adhiana46/aegis-asesmen/internal/user/data/source"
	PkgDataSource "github.com/Adhiana46/aegis-asesmen/pkg/data_sources"
	"github.com/pkg/errors"
)

type userPersistentPostgres struct {
	db *PkgDataSource.PostgresDB
}

func NewUserPersistentPostgres(db *PkgDataSource.PostgresDB) Source.IUserPersistent {
	return &userPersistentPostgres{
		db: db,
	}
}

func (r *userPersistentPostgres) GetByID(ctx context.Context, id string) (*Model.User, error) {
	query := `
		SELECT
			id,
			email,
			password,
			role,
			created_at,
			updated_at
		FROM "users"
		WHERE id = $1
	`

	model := Model.User{}
	err := r.db.GetContext(ctx, &model, query, id)
	if err == nil {
		return &model, nil
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return nil, errors.Wrap(err, "error GetByID")
}

func (r *userPersistentPostgres) GetByEmail(ctx context.Context, email string) (*Model.User, error) {
	query := `
		SELECT
			id,
			email,
			password,
			role,
			created_at,
			updated_at
		FROM "users"
		WHERE email = $1
	`

	model := Model.User{}
	err := r.db.GetContext(ctx, &model, query, email)
	if err == nil {
		return &model, nil
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return nil, errors.Wrap(err, "error GetByEmail")
}

func (r *userPersistentPostgres) Store(ctx context.Context, userModel *Model.User) error {
	query := `
		INSERT INTO "users"
		(id, email, password, role, created_at, updated_at)
		VALUES
		(:id, :email, :password, :role, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, userModel)
	if err != nil {
		return errors.Wrap(err, "error Store")
	}

	return nil
}

func (r *userPersistentPostgres) Update(ctx context.Context, userModel *Model.User) error {
	query := `
		UPDATE "users"
			SET
				email = :email,
				password = :password,
				role = :role,
				created_at = :created_at,
				updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, userModel)
	if err != nil {
		return errors.Wrap(err, "error Update")
	}

	return nil
}

func (r *userPersistentPostgres) Destroy(ctx context.Context, userModel *Model.User) error {
	query := `
		DELETE FROM "users" WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userModel.Id)
	if err != nil {
		return errors.Wrap(err, "error Destroy")
	}

	return nil
}
