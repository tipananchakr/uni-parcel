package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tipananchakr/uni-parcel/internals/application"
	"github.com/tipananchakr/uni-parcel/internals/core/domain"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AuthHandler struct {
	authService *application.AuthService
}

type authResponse struct {
	User  domain.User `json:"user"`
	Token string      `json:"token"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterAuthRoutes(route fiber.Router, authService *application.AuthService) {
	handler := &AuthHandler{
		authService: authService,
	}

	route.Post("/register", handler.Register)
	route.Post("/login", handler.Login)
}

func (a *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	result, err := a.authService.Register(c.Context(), domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(authResponse{
		User:  result.User,
		Token: result.Token,
	})
}

func (a *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	result, err := a.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(authResponse{
		User:  result.User,
		Token: result.Token,
	})
}
