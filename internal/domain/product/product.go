package product

import (
	"errors"
	"time"
)

// ProductID is a Value Object for product identification.
type ProductID string

// Product is an Entity/Aggregate Root.
type Product struct {
	ID          ProductID
	Name        string
	Description string
	Price       float64
	SKU         string // Stock Keeping Unit
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct creates a new Product entity.
func NewProduct(id ProductID, name, description, sku string, price float64) (*Product, error) {
	if name == "" || sku == "" || price <= 0 {
		return nil, errors.New("product name, SKU, and positive price are required")
	}
	return &Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		SKU:         sku,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// UpdateDetails updates product name, description, and price.
func (p *Product) UpdateDetails(newName, newDescription string, newPrice float64) error {
	if newName == "" || newPrice <= 0 {
		return errors.New("product name and positive price are required for update")
	}
	p.Name = newName
	p.Description = newDescription
	p.Price = newPrice
	p.UpdatedAt = time.Now()
	return nil
}
