package apis

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/database/entities"
	"gorm.io/gorm"
)

type TestApi struct {
	db *gorm.DB
}

func NewTestApi(dbInstance *gorm.DB) *TestApi {
	return &TestApi{db: dbInstance}
}

func (api *TestApi) DbTest(c *gin.Context) {
	var user entities.User
	api.db.First(&user)

	task := entities.Task{
		Title:       "AAAAAAAAAA",
		Description: user.Nickname,
		State:       func(i int) *int { return &i }(1),
		User:        user,
		Label: []*entities.Label{{
			Name: "Ahoy",
		}},
		StartTime: func(i time.Time) *time.Time { return &i }(time.Now()),
		EndTime:   func(i time.Time) *time.Time { return &i }(time.Now().AddDate(0, 0, 1)),
	}
	api.db.Create(&task)
	c.JSON(http.StatusOK, task.ID)
}

func (api *TestApi) Repeat(c *gin.Context) {
	repeat, err := strconv.Atoi(os.Getenv("APPSETTING_REPEAT"))
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go!\n")
	}
	c.String(http.StatusOK, buffer.String())
}
