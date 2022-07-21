package apis

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/auth"
	"github.com/pc01pc013/task-management/database/entities"
	modelsAuth "github.com/pc01pc013/task-management/models/auth"
	"github.com/pc01pc013/task-management/utils"
	"gorm.io/gorm"
)

type AuthApi struct {
	db *gorm.DB
}

func NewAuthApi(dbInstance *gorm.DB) *AuthApi {
	return &AuthApi{db: dbInstance}
}

func (api *AuthApi) Login(c *gin.Context) {
	var req modelsAuth.LoginReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var userEntities entities.User
	if result := api.db.Where("Username = ?", req.Username).First(&userEntities); result.Error != nil {
		log.Printf("Login Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("Login Failed. 1")
		c.JSON(http.StatusOK, context)
		return
	}

	if strings.Compare(*userEntities.Password, *req.Password) != 0 {
		log.Printf("Login Error: Password not equal.")
		context := utils.MakeResponseResultFailed("Login Failed. 2")
		c.JSON(http.StatusOK, context)
		return
	}

	ss, err := auth.NewToken(*req.Username)
	log.Printf("%v %v", ss, err)

	context := utils.MakeResponseResultSuccess(ss)
	c.JSON(http.StatusOK, context)
}
