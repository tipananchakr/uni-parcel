package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tipananchakr/uni-parcel/internals/application"
	"github.com/tipananchakr/uni-parcel/internals/core/domain"
)

type MajorHandler struct {
	majorService *application.MajorService
	authService  *application.AuthService
}

func RegisterMajorRoutes(group fiber.Router, majorService *application.MajorService, authService *application.AuthService) {
	handler := MajorHandler{
		majorService: majorService,
		authService:  authService,
	}

	group.Get("/", handler.GetAll)
	group.Get("/:id", handler.GetByID)
	group.Post("/", handler.Create)
	group.Patch("/:id", handler.Update)
	group.Delete("/:id", handler.Delete)
}

type majorResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type majorRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (h *MajorHandler) GetAll(c *fiber.Ctx) error {
	majors, err := h.majorService.GetAllMajors(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch majors: " + err.Error()})
	}

	majorResponses := make([]majorResponse, len(majors))
	for i, major := range majors {
		majorResponses[i] = majorResponse{
			ID:   major.ID.Hex(),
			Code: major.Code,
			Name: major.Name,
		}
	}

	return c.JSON(majorResponses)
}

func (h *MajorHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	major, err := h.majorService.GetMajorByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Major not found"})
	}

	return c.JSON(majorResponse{
		ID:   major.ID.Hex(),
		Code: major.Code,
		Name: major.Name,
	})

}

func (h *MajorHandler) Create(c *fiber.Ctx) error {
	var req majorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
	}

	major := &domain.Major{
		Code: req.Code,
		Name: req.Name,
	}

	err := h.majorService.CreateMajor(c.Context(), major)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create major: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Major created successfully"})
}

func (h *MajorHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req domain.MajorUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
	}

	if req.Code == nil && req.Name == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one field (code or name) must be provided for update"})
	}

	err := h.majorService.UpdateMajor(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update major: " + err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "major updated successfully"})
}

func (h *MajorHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.majorService.DeleteMajor(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete major: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "major deleted successfully"})
}
