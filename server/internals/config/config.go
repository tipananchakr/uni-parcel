package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	DatabaseName   string
	JwtSecret      string
	UserCollection string
	DormCollection string
}

func Load() (Config, error) {
	godotenv.Load(".env")

	cfg := Config{
		Port:           getEnv("PORT", "3000"),
		MongoURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName:   getEnv("DATABASE_NAME", "uni_parcel"),
		JwtSecret:      getEnv("JWT_SECRET", "your_jwt_secret"),
		UserCollection: getEnv("USER_COLLECTION", "users"),
		DormCollection: getEnv("DORM_COLLECTION", "dorms"),
	}
	return cfg, nil
}

func getEnv(name string, fallback string) string {
	value := os.Getenv(name)

	if value == "" {
		return fallback
	}
	return value
}
