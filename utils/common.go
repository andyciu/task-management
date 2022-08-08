package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/database/entities"
	"github.com/pc01pc013/task-management/enums/response"
	"github.com/pc01pc013/task-management/models"
	"gorm.io/gorm"
)

func MakeResponseResultSuccess(content interface{}) models.CommonRes {
	return models.CommonRes{
		Code:    string(response.Success),
		Message: "",
		Content: content,
	}
}

func MakeResponseResult(code response.CommonResCode, message string, content interface{}) models.CommonRes {
	return models.CommonRes{
		Code:    string(code),
		Message: message,
		Content: content,
	}
}

func MakeResponseResultFailed(message string) models.CommonRes {
	return models.CommonRes{
		Code:    string(response.Failure),
		Message: message,
		Content: nil,
	}
}

func JsonNumberPointToIntPoint(jnum *json.Number) *int {
	if jnum == nil {
		return nil
	}

	tempnum, _ := jnum.Int64()
	result := int(tempnum)

	return &result
}

func QueryCondLikeString(str string) string {
	return "%" + str + "%"
}

func DateTimePrase(value *time.Time) *time.Time {
	if value != nil {
		newtime, _ := time.Parse("2006/01/02", value.Format("2006/01/02"))
		return &newtime
	} else {
		return nil
	}
}

func GetUserID(c *gin.Context, dbInstance *gorm.DB) (uint, error) {
	username, _ := c.Get("username")
	authtype, _ := c.Get("authtype")
	if username == nil || authtype == nil {
		return 0, fmt.Errorf("gin.Context find error")
	}
	var userEntities entities.User
	if result := dbInstance.Where("Username = ? AND auth_type = ?", username, authtype).Select("ID").First(&userEntities); result.Error != nil {
		return 0, result.Error
	}

	return userEntities.ID, nil
}
