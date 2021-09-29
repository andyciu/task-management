package apis

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LabelsApi struct {
	db *gorm.DB
}

func NewLabelsApi(dbInstance *gorm.DB) *LabelsApi {
	return &LabelsApi{db: dbInstance}
}

func (api *LabelsApi) List(c *gin.Context) {

}
