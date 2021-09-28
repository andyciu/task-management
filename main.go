package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	"github.com/pc01pc013/task-management/database"
	"github.com/pc01pc013/task-management/models"
)

func repeatHandler(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/repeat", repeatHandler(repeat))

	router.GET("/db", dbTest(database.GetDBInstance()))

	router.Run(":" + port)
}

func dbTest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		db.First(&user)

		task := models.Task{
			Title:       "AAAAAAAAAA",
			Description: user.Nickname,
			State:       func(i int) *int { return &i }(1),
			User:        user,
			Label: []*models.Label{{
				Name: "Ahoy",
			}},
			StartTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 0, 1),
		}
		db.Create(&task)
	}
}
