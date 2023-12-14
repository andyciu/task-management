package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/pc01pc013/task-management/database"
	"github.com/pc01pc013/task-management/router"
)

func main() {
	defer database.Close()

	// Azure App Service sets the port as an Environment Variable
	// This can be random, so needs to be loaded at startup
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router := router.InitRouter()
	router.Run("127.0.0.1:" + port)
}
