package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/pc01pc013/task-management/database"
	"github.com/pc01pc013/task-management/router"
)

func main() {
	defer database.Close()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := router.InitRouter()
	router.Run(":" + port)
}
