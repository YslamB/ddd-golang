package postgres

import (
	"context"
	"database/sql"
	"fmt"

	domain_common "gddd/internal/domain/common"
	domain_user "gddd/internal/domain/user" // Import domain package

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresUserRepository implements user.Repository for PostgreSQL.
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserRepository creates a new Postgres user repository.
func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Save implements user.Repository.Save. It handles inserting new users or updating existing ones.
func (r *PostgresUserRepository) Save(ctx context.Context, u *domain_user.User) error {
	query := `INSERT INTO users (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)
	          ON CONFLICT (id) DO UPDATE SET email = $2, password = $3, updated_at = $5`
	_, err := r.db.Exec(ctx, query, u.ID, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		// Wrap the database error to provide context.
		return fmt.Errorf("postgres: failed to save user %s: %w", u.ID, err)
	}
	return nil
}

// FindByID implements user.Repository.FindByID. It retrieves a user by their unique ID.
func (r *PostgresUserRepository) FindByID(ctx context.Context, id domain_user.UserID) (*domain_user.User, error) {
	u := &domain_user.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Map sql.ErrNoRows to a domain-specific NotFoundError.
			return nil, domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found", id))
		}
		// Wrap other database errors.
		return nil, fmt.Errorf("postgres: failed to find user by ID %s: %w", id, err)
	}
	return u, nil
}

// FindByEmail implements user.Repository.FindByEmail. It retrieves a user by their email address.
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain_user.User, error) {
	u := &domain_user.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Map sql.ErrNoRows to a domain-specific NotFoundError.
			return nil, domain_common.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
		}
		// Wrap other database errors.
		return nil, fmt.Errorf("postgres: failed to find user by email %s: %w", email, err)
	}
	return u, nil
}

// Delete implements user.Repository.Delete. It removes a user from the database.
func (r *PostgresUserRepository) Delete(ctx context.Context, id domain_user.UserID) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres: failed to delete user %s: %w", id, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		// If no rows were affected, it means the user wasn't found for deletion.
		return domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found for deletion", id))
	}
	return nil
}
