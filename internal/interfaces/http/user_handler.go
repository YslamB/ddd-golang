package http

import (
	app_user "gddd/internal/application/user"
	"gddd/internal/application/user/dtos"
	domain_common "gddd/internal/domain/common"
	infra_logging "gddd/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *app_user.Service
	logger      *infra_logging.Logger
}

func NewUserHandler(userService *app_user.Service, logger *infra_logging.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	h.logger.Infof("Received RegisterUser request from %s", c.IP())

	var req dtos.RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warnf("Invalid request body for RegisterUser: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userResp, err := h.userService.RegisterUser(c.Context(), req)
	if err != nil {
		h.logger.Errorf("Error registering user: %v", err)
		switch {
		case domain_common.IsValidationError(err):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case domain_common.IsConflictError(err):
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	h.logger.Infof("User %s registered successfully.", userResp.ID)
	return c.Status(fiber.StatusCreated).JSON(userResp)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	h.logger.Infof("Received GetUserByID request for ID: %s from %s", userID, c.IP())

	userResp, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		h.logger.Errorf("Error getting user %s: %v", userID, err)
		switch {
		case domain_common.IsNotFoundError(err):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	h.logger.Infof("User %s retrieved successfully.", userResp.ID)
	return c.Status(fiber.StatusOK).JSON(userResp)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	h.logger.Infof("Received UpdateUser request for ID: %s from %s", userID, c.IP())

	var req dtos.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warnf("Invalid request body for UpdateUser %s: %v", userID, err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userResp, err := h.userService.UpdateUser(c.Context(), userID, req)
	if err != nil {
		h.logger.Errorf("Error updating user %s: %v", userID, err)
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

	h.logger.Infof("User %s updated successfully.", userResp.ID)
	return c.Status(fiber.StatusOK).JSON(userResp)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	h.logger.Infof("Received DeleteUser request for ID: %s from %s", userID, c.IP())

	err := h.userService.DeleteUser(c.Context(), userID)
	if err != nil {
		h.logger.Errorf("Error deleting user %s: %v", userID, err)
		switch {
		case domain_common.IsNotFoundError(err):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	h.logger.Infof("User %s deleted successfully.", userID)
	return c.Status(fiber.StatusNoContent).SendString("")
}
