package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
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

// func dbFunc(db *sql.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
// 			c.String(http.StatusInternalServerError,
// 				fmt.Sprintf("Error creating database table: %q", err))
// 			return
// 		}

// 		if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
// 			c.String(http.StatusInternalServerError,
// 				fmt.Sprintf("Error incrementing tick: %q", err))
// 			return
// 		}

// 		rows, err := db.Query("SELECT tick FROM ticks")
// 		if err != nil {
// 			c.String(http.StatusInternalServerError,
// 				fmt.Sprintf("Error reading ticks: %q", err))
// 			return
// 		}

// 		defer rows.Close()
// 		for rows.Next() {
// 			var tick time.Time
// 			if err := rows.Scan(&tick); err != nil {
// 				c.String(http.StatusInternalServerError,
// 					fmt.Sprintf("Error scanning ticks: %q", err))
// 				return
// 			}
// 			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
// 		}
// 	}
// }

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

	sqlDB, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database (gorm): %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/repeat", repeatHandler(repeat))

	router.GET("/db", dbTest(gormDB))

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
