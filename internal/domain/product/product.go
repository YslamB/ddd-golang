package product

import (
	"errors"
	"time"
)

type ProductID string

type Product struct {
	ID          ProductID
	Name        string
	Description string
	Price       float64
	SKU         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

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
