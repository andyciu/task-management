package apis

import (
	"log"
	"net/http"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/database/entities"
	"github.com/pc01pc013/task-management/models/labels"
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
	var labelEntities []entities.Label
	api.db.Find(&labelEntities)

	var result []labels.LabelListRes
	From(labelEntities).
		Select(func(i interface{}) interface{} {
			return labels.LabelListRes{
				ID:   int(i.(entities.Label).ID),
				Name: i.(entities.Label).Name,
			}
		}).ToSlice(&result)

	context := utils.MakeResponseResultSuccess(result)

	c.JSON(http.StatusOK, context)
}

func (api *LabelsApi) Create(c *gin.Context) {
	var req labels.LabelCreateReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	newlabel := entities.Label{
		Name: req.Name,
	}
	api.db.Create(&newlabel)

	context := utils.MakeResponseResultSuccess(newlabel.ID)
	c.JSON(http.StatusOK, context)
}
