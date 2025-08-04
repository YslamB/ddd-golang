package http

import (
	app_product "gddd/internal/application/product"
	"gddd/internal/application/product/dtos" // Import DTOs
	domain_common "gddd/internal/domain/common"
	infra_logging "gddd/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
)

// ProductHandler handles HTTP requests for product-related operations.
type ProductHandler struct {
	productService *app_product.Service
	logger         *infra_logging.Logger
}

// NewProductHandler creates a new ProductHandler instance.
func NewProductHandler(productService *app_product.Service, logger *infra_logging.Logger) *ProductHandler {
	return &ProductHandler{productService: productService, logger: logger}
}

// CreateProduct handles POST /api/v1/products requests.
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	h.logger.Infof("Received CreateProduct request from %s", c.IP())

	var req dtos.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warnf("Invalid request body for CreateProduct: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	productResp, err := h.productService.CreateProduct(c.Context(), req)
	if err != nil {
		h.logger.Errorf("Error creating product: %v", err)
		switch {
		case domain_common.IsValidationError(err):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case domain_common.IsConflictError(err):
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	h.logger.Infof("Product %s created successfully.", productResp.ID)
	return c.Status(fiber.StatusCreated).JSON(productResp)
}

// GetProductByID handles GET /api/v1/products/:id requests.
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	productID := c.Params("id")
	h.logger.Infof("Received GetProductByID request for ID: %s from %s", productID, c.IP())

	productResp, err := h.productService.GetProductByID(c.Context(), productID)
	if err != nil {
		h.logger.Errorf("Error getting product %s: %v", productID, err)
		switch {
		case domain_common.IsNotFoundError(err):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	h.logger.Infof("Product %s retrieved successfully.", productResp.ID)
	return c.Status(fiber.StatusOK).JSON(productResp)
}

// UpdateProduct handles PUT /api/v1/products/:id requests.
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	h.logger.Infof("Received UpdateProduct request for ID: %s from %s", productID, c.IP())

	var req dtos.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warnf("Invalid request body for UpdateProduct %s: %v", productID, err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	productResp, err := h.productService.UpdateProduct(c.Context(), productID, req)
	if err != nil {
		h.logger.Errorf("Error updating product %s: %v", productID, err)
		switch {
		case domain_common.IsValidationError(err):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case domain_common.IsNotFoundError(err):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case domain_common.IsConflictError(err):
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	h.logger.Infof("Product %s updated successfully.", productResp.ID)
	return c.Status(fiber.StatusOK).JSON(productResp)
}

// DeleteProduct handles DELETE /api/v1/products/:id requests.
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	h.logger.Infof("Received DeleteProduct request for ID: %s from %s", productID, c.IP())

	err := h.productService.DeleteProduct(c.Context(), productID)
	if err != nil {
		h.logger.Errorf("Error deleting product %s: %v", productID, err)
		switch {
		case domain_common.IsNotFoundError(err):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	h.logger.Infof("Product %s deleted successfully.", productID)
	return c.Status(fiber.StatusNoContent).SendString("") // No content for successful deletion
}
