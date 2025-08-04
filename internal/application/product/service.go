package product

import (
	"context"
	"fmt"

	"gddd/internal/application/product/dtos" // Import DTOs
	domain_common "gddd/internal/domain/common"
	domain_product "gddd/internal/domain/product" // Import domain package
	infra_logging "gddd/internal/infrastructure/logging"
	shared_validation "gddd/internal/shared/validation"

	"github.com/google/uuid"
)

// ProductRepository defines the interface for the product repository,
// used by the application service.
type ProductRepository interface {
	Save(ctx context.Context, product *domain_product.Product) error
	FindByID(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error)
	FindBySKU(ctx context.Context, sku string) (*domain_product.Product, error)
	Delete(ctx context.Context, id domain_product.ProductID) error
}

// Service defines the application service for product-related operations.
type Service struct {
	productRepo ProductRepository
	logger      *infra_logging.Logger
}

// NewService creates a new product application service.
func NewService(productRepo ProductRepository, logger *infra_logging.Logger) *Service {
	return &Service{productRepo: productRepo, logger: logger}
}

// CreateProduct handles the product creation use case.
func (s *Service) CreateProduct(ctx context.Context, req dtos.CreateProductRequest) (*dtos.ProductResponse, error) {
	s.logger.Infof("Attempting to create product with SKU: %s", req.SKU)

	if err := shared_validation.ValidateStruct(req); err != nil {
		s.logger.Warnf("Validation error for product creation: %v", err)
		return nil, domain_common.NewValidationError(err.Error())
	}

	existingProduct, err := s.productRepo.FindBySKU(ctx, req.SKU)
	if err != nil && !domain_common.IsNotFoundError(err) {
		s.logger.Errorf("Failed to check existing product by SKU %s: %v", req.SKU, err)
		return nil, fmt.Errorf("failed to check existing product: %w", err)
	}
	if existingProduct != nil {
		s.logger.Warnf("Product with SKU %s already exists.", req.SKU)
		return nil, domain_common.NewConflictError("product with this SKU already exists")
	}

	newProductID := domain_product.ProductID(uuid.New().String())
	newProduct, err := domain_product.NewProduct(newProductID, req.Name, req.Description, req.SKU, req.Price)
	if err != nil {
		s.logger.Errorf("Failed to create new product domain entity for SKU %s: %v", req.SKU, err)
		return nil, fmt.Errorf("failed to create new product domain entity: %w", err)
	}

	if err := s.productRepo.Save(ctx, newProduct); err != nil {
		s.logger.Errorf("Failed to save product %s: %v", newProduct.ID, err)
		return nil, fmt.Errorf("failed to save product: %w", err)
	}

	s.logger.Infof("Product %s created successfully.", newProduct.ID)
	return &dtos.ProductResponse{
		ID:          string(newProduct.ID),
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		SKU:         newProduct.SKU,
		CreatedAt:   newProduct.CreatedAt,
		UpdatedAt:   newProduct.UpdatedAt,
	}, nil
}

// GetProductByID handles retrieving a product by ID use case.
func (s *Service) GetProductByID(ctx context.Context, id string) (*dtos.ProductResponse, error) {
	s.logger.Infof("Attempting to get product by ID: %s", id)

	productID := domain_product.ProductID(id)
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		if domain_common.IsNotFoundError(err) {
			s.logger.Warnf("Product with ID %s not found.", id)
			return nil, domain_common.NewNotFoundError(fmt.Sprintf("product with ID %s not found", id))
		}
		s.logger.Errorf("Failed to get product by ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	s.logger.Infof("Product %s retrieved successfully.", product.ID)
	return &dtos.ProductResponse{
		ID:          string(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

// UpdateProduct handles updating a product's details use case.
func (s *Service) UpdateProduct(ctx context.Context, id string, req dtos.UpdateProductRequest) (*dtos.ProductResponse, error) {
	s.logger.Infof("Attempting to update product %s", id)

	if err := shared_validation.ValidateStruct(req); err != nil {
		s.logger.Warnf("Validation error for product update %s: %v", id, err)
		return nil, domain_common.NewValidationError(err.Error())
	}

	productID := domain_product.ProductID(id)
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		if domain_common.IsNotFoundError(err) {
			s.logger.Warnf("Product with ID %s not found for update.", id)
			return nil, domain_common.NewNotFoundError(fmt.Sprintf("product with ID %s not found", id))
		}
		s.logger.Errorf("Failed to find product %s for update: %v", id, err)
		return nil, fmt.Errorf("failed to retrieve product for update: %w", err)
	}

	if err := product.UpdateDetails(req.Name, req.Description, req.Price); err != nil {
		s.logger.Errorf("Failed to update product %s details: %v", product.ID, err)
		return nil, fmt.Errorf("failed to update product details: %w", err)
	}

	if err := s.productRepo.Save(ctx, product); err != nil {
		s.logger.Errorf("Failed to save updated product %s: %v", product.ID, err)
		return nil, fmt.Errorf("failed to save updated product: %w", err)
	}

	s.logger.Infof("Product %s updated successfully.", product.ID)
	return &dtos.ProductResponse{
		ID:          string(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

// DeleteProduct handles deleting a product use case.
func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	s.logger.Infof("Attempting to delete product: %s", id)

	productID := domain_product.ProductID(id)
	err := s.productRepo.Delete(ctx, productID)
	if err != nil {
		if domain_common.IsNotFoundError(err) {
			s.logger.Warnf("Product with ID %s not found for deletion.", id)
			return domain_common.NewNotFoundError(fmt.Sprintf("product with ID %s not found", id))
		}
		s.logger.Errorf("Failed to delete product %s: %v", id, err)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	s.logger.Infof("Product %s deleted successfully.", id)
	return nil
}
