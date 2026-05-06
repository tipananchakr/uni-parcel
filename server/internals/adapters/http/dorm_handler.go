package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tipananchakr/uni-parcel/internals/application"
	"github.com/tipananchakr/uni-parcel/internals/core/domain"
)

type DormHandler struct {
	dormService *application.DormService
	authService *application.AuthService
}

type dormRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type dormResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func RegisDormRoute(router fiber.Router, dormService *application.DormService, authService *application.AuthService) {
	handler := DormHandler{
		dormService: dormService,
		authService: authService,
	}

	router.Get("/", handler.GetAll)
	router.Get("/:id", handler.GetByID)
	router.Post("/", handler.Create)
	router.Patch("/:id", handler.Update)
	router.Delete("/:id", handler.Delete)
}

func (d *DormHandler) GetAll(c *fiber.Ctx) error {
	dorms, err := d.dormService.GetAllDorms(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch dorms: " + err.Error()})
	}

	dormResponses := make([]dormResponse, len(dorms))
	for i, dorm := range dorms {
		dormResponses[i] = dormResponse{
			ID:   dorm.ID.Hex(),
			Code: dorm.Code,
			Name: dorm.Name,
		}
	}

	return c.JSON(dormResponses)
}

func (d *DormHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	dorm, err := d.dormService.GetDormByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Dorm not found"})
	}

	dormResponse := dormResponse{
		ID:   dorm.ID.Hex(),
		Code: dorm.Code,
		Name: dorm.Name,
	}

	return c.JSON(dormResponse)
}

func (d *DormHandler) Create(c *fiber.Ctx) error {
	var req dormRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.Code == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "code and name are required"})
	}

	dorm := &domain.Dorm{
		Code: req.Code,
		Name: req.Name,
	}

	err := d.dormService.CreateDorm(c.Context(), dorm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create dorm: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(dorm)
}

func (d *DormHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req domain.DormUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Code == nil && req.Name == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no fields to update"})
	}

	err := d.dormService.UpdateDorm(c.Context(), id, domain.DormUpdate{
		Code: req.Code,
		Name: req.Name,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update dorm: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Dorm updated successfully"})

}

func (d *DormHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := d.dormService.DeleteDorm(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete dorm: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Dorm deleted successfully"})
}
