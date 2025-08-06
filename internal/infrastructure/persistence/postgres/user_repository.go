package postgres

import (
	"context"
	"database/sql"
	"fmt"

	domain_common "gddd/internal/domain/common"
	domain_user "gddd/internal/domain/user"
	db "gddd/pkg/database"
)

type PostgresUserRepository struct {
	db *db.Database
}

func NewPostgresUserRepository(db *db.Database) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, u *domain_user.User) error {
	query := `INSERT INTO users (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)
	          ON CONFLICT (id) DO UPDATE SET email = $2, password = $3, updated_at = $5`
	_, err := r.db.Exec(ctx, query, u.ID, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {

		return fmt.Errorf("postgres: failed to save user %s: %w", u.ID, err)
	}
	return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id domain_user.UserID) (*domain_user.User, error) {
	u := &domain_user.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found", id))
		}

		return nil, fmt.Errorf("postgres: failed to find user by ID %s: %w", id, err)
	}
	return u, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain_user.User, error) {
	u := &domain_user.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, domain_common.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
		}

		return nil, fmt.Errorf("postgres: failed to find user by email %s: %w", email, err)
	}
	return u, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id domain_user.UserID) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres: failed to delete user %s: %w", id, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {

		return domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found for deletion", id))
	}
	return nil
}
