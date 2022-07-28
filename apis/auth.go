package apis

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/auth"
	"github.com/pc01pc013/task-management/database/entities"
	modelsAuth "github.com/pc01pc013/task-management/models/auth"
	"github.com/pc01pc013/task-management/utils"
	"golang.org/x/crypto/scrypt"
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
		context := utils.MakeResponseResultFailed("Login Failed.")
		c.JSON(http.StatusOK, context)
		return
	}

	dk, err := scrypt.Key([]byte(*req.Password), []byte(os.Getenv("PASSWORDSALT")), 32768, 8, 1, 32)
	if err != nil {
		log.Printf("Login Error: %q", err)
		context := utils.MakeResponseResultFailed("Login Failed.")
		c.JSON(http.StatusOK, context)
		return
	}
	if strings.Compare(*userEntities.Password, base64.StdEncoding.EncodeToString(dk)) != 0 {
		log.Printf("Login Error: Password not equal.")
		log.Printf(base64.StdEncoding.EncodeToString(dk))
		context := utils.MakeResponseResultFailed("Login Failed.")
		c.JSON(http.StatusOK, context)
		return
	}

	ss, err := auth.NewToken(*req.Username)
	log.Printf("%v %v", ss, err)

	context := utils.MakeResponseResultSuccess(ss)
	c.JSON(http.StatusOK, context)
}
