package inmemory

import (
	"context"
	"fmt"
	"sync"

	domain_common "gddd/internal/domain/common"
	domain_product "gddd/internal/domain/product"
)

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[domain_product.ProductID]*domain_product.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[domain_product.ProductID]*domain_product.Product),
	}
}

func (r *InMemoryProductRepository) Save(ctx context.Context, p *domain_product.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[p.ID] = p
	return nil
}

func (r *InMemoryProductRepository) FindByID(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, ok := r.products[id]
	if !ok {
		return nil, domain_common.NewNotFoundError(fmt.Sprintf("product with ID %s not found", id))
	}
	copiedProduct := *product
	return &copiedProduct, nil
}

func (r *InMemoryProductRepository) FindBySKU(ctx context.Context, sku string) (*domain_product.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, product := range r.products {
		if product.SKU == sku {
			copiedProduct := *product
			return &copiedProduct, nil
		}
	}
	return nil, domain_common.NewNotFoundError(fmt.Sprintf("product with SKU %s not found", sku))
}

func (r *InMemoryProductRepository) Delete(ctx context.Context, id domain_product.ProductID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[id]; !ok {
		return domain_common.NewNotFoundError(fmt.Sprintf("product with ID %s not found for deletion", id))
	}
	delete(r.products, id)
	return nil
}
