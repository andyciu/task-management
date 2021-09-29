package apis

import (
	"net/http"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/database/entities"
	modelsRes "github.com/pc01pc013/task-management/models/labels"
	"github.com/pc01pc013/task-management/utils"
	"gorm.io/gorm"
)

type LabelsApi struct {
	db *gorm.DB
}

func NewLabelsApi(dbInstance *gorm.DB) *LabelsApi {
	return &LabelsApi{db: dbInstance}
}

func (api *LabelsApi) List(c *gin.Context) {
	var labels []entities.Label
	api.db.Find(&labels)

	var result []modelsRes.LabelListRes
	From(&labels).
		Select(func(i interface{}) interface{} {
			return modelsRes.LabelListRes{
				ID:   int(i.(entities.Label).ID),
				Name: i.(entities.Label).Name,
			}
		}).ToSlice(&result)

	context := utils.MakeResponseResultSuccess(result)

	c.JSON(http.StatusOK, context)
}
