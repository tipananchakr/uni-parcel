package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	mongorepo "github.com/tipananchakr/uni-parcel/internals/adapters/mongodb"
	"github.com/tipananchakr/uni-parcel/internals/adapters/security"
	"github.com/tipananchakr/uni-parcel/internals/config"

	"github.com/tipananchakr/uni-parcel/internals/adapters/http"
	"github.com/tipananchakr/uni-parcel/internals/application"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize MongoDB
	ctx := context.Background()
	mongoClient, err := mongorepo.Connect(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)

	database := mongoClient.Database(cfg.DatabaseName)

	// Auth
	userRepository := mongorepo.NewUserRepository(database.Collection(cfg.UserCollection))
	hasher := security.NewBcryptPasswordHasher()
	tokenManager := security.NewHMACTokenManager(cfg.JwtSecret, 24*time.Hour)
	authService := application.NewAuthService(userRepository, hasher, tokenManager)

	// Dorm
	dormRepository := mongorepo.NewDormRepository(database.Collection(cfg.DormCollection))
	dormService := application.NewDormService(ctx, dormRepository)

	// Initialize Fiber app and register routes
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	http.RegisterRoutes(app, http.Services{
		Auth: authService,
		Dorm: dormService,
	})

	log.Fatal(app.Listen("0.0.0.0:" + cfg.Port))
}
