package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tipananchakr/uni-parcel/internals/application"
	"github.com/tipananchakr/uni-parcel/internals/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authService *application.AuthService
}

type UserResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
	Role  string             `json:"role"`
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
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
		Role:     domain.RoleUser,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(authResponse{
		User: UserResponse{
			ID:    result.User.ID,
			Name:  result.User.Name,
			Email: result.User.Email,
			Role:  string(result.User.Role),
		},
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
		User: UserResponse{
			ID:    result.User.ID,
			Name:  result.User.Name,
			Email: result.User.Email,
			Role:  string(result.User.Role),
		},
		Token: result.Token,
	})
}
