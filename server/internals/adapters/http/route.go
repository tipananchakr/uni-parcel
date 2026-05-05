package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tipananchakr/uni-parcel/internals/application"
)

type Services struct {
	Auth *application.AuthService
	Dorm *application.DormService
}

func RegisterRoutes(app *fiber.App, services Services) {
	api := app.Group("/api")
	RegisterAuthRoutes(api.Group("/auth"), services.Auth)
	RegisDormRoute(api.Group("/dorms"), services.Dorm, services.Auth)
}
