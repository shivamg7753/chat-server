package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port        string
	DBPath      string
	JWTSecret   string
	JWTExpiry   string
)

func Load() {
	godotenv.Load()
	
	Port = getEnv("PORT", "3000")
	DBPath = getEnv("DB_PATH", "./chat.db")
	JWTSecret = getEnv("JWT_SECRET", "default-secret-key")
	JWTExpiry = getEnv("JWT_EXPIRY_HOURS", "24")
	
	log.Println("Config loaded")
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
