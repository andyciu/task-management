package apis

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/database/entities"
	"github.com/pc01pc013/task-management/models/user"
	"github.com/pc01pc013/task-management/utils"
	"gorm.io/gorm"
)

type UserApi struct {
	db *gorm.DB
}

func NewUserApi(dbInstance *gorm.DB) *UserApi {
	return &UserApi{db: dbInstance}
}

func (api *UserApi) GetNickName(c *gin.Context) {
	var userEntities *entities.User
	usernamestr, _ := c.Get("username")
	tx := api.db.Where("Username = ?", usernamestr)

	if result := tx.Find(&userEntities); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	result := user.UserGetNickNameRes{
		Name: *userEntities.Nickname,
	}

	context := utils.MakeResponseResultSuccess(result)
	c.JSON(http.StatusOK, context)
}
