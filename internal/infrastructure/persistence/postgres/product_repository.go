package postgres

import (
	"context"
	"database/sql"
	"fmt"

	domain_common "gddd/internal/domain/common"
	domain_product "gddd/internal/domain/product"
)

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Save(ctx context.Context, p *domain_product.Product) error {
	query := `INSERT INTO products (id, name, description, price, sku, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)
	          ON CONFLICT (id) DO UPDATE SET name = $2, description = $3, price = $4, sku = $5, updated_at = $7`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Description, p.Price, p.SKU, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres: failed to save product %s: %w", p.ID, err)
	}
	return nil
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
	p := &domain_product.Product{}
	query := `SELECT id, name, description, price, sku, created_at, updated_at FROM products WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.SKU, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_common.NewNotFoundError(fmt.Sprintf("product with ID %s not found", id))
		}
		return nil, fmt.Errorf("postgres: failed to find product by ID %s: %w", id, err)
	}
	return p, nil
}

func (r *PostgresProductRepository) FindBySKU(ctx context.Context, sku string) (*domain_product.Product, error) {
	p := &domain_product.Product{}
	query := `SELECT id, name, description, price, sku, created_at, updated_at FROM products WHERE sku = $1`
	err := r.db.QueryRowContext(ctx, query, sku).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.SKU, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain_common.NewNotFoundError(fmt.Sprintf("product with SKU %s not found", sku))
		}
		return nil, fmt.Errorf("postgres: failed to find product by SKU %s: %w", sku, err)
	}
	return p, nil
}

func (r *PostgresProductRepository) Delete(ctx context.Context, id domain_product.ProductID) error {
	query := `DELETE FROM products WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres: failed to delete product %s: %w", id, err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain_common.NewNotFoundError(fmt.Sprintf("product with ID %s not found for deletion", id))
	}
	return nil
}
