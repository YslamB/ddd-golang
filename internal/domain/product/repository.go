package product

import "context"

// Repository defines the interface for Product persistence operations.
type Repository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id ProductID) (*Product, error)
	FindBySKU(ctx context.Context, sku string) (*Product, error)
	Delete(ctx context.Context, id ProductID) error
}
