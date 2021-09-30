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
	if result := api.db.Create(&newlabel); result.Error != nil {
		log.Printf("Create Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context := utils.MakeResponseResultSuccess(newlabel.ID)
	c.JSON(http.StatusOK, context)
}

func (api *LabelsApi) Update(c *gin.Context) {
	var req labels.LabelUpdateReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var labelEntities entities.Label
	reqID, _ := req.ID.Int64()
	if result := api.db.First(&labelEntities, reqID); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if labelEntities.Name != req.Name {
		labelEntities.Name = req.Name
	}

	if result := api.db.Save(&labelEntities); result.Error != nil {
		log.Printf("Save Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context := utils.MakeResponseResultSuccess(nil)
	c.JSON(http.StatusOK, context)
}

func (api *LabelsApi) Delete(c *gin.Context) {
	var req labels.LabelDeleteReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var labelEntities entities.Label
	reqID, _ := req.ID.Int64()
	if result := api.db.First(&labelEntities, reqID); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if result := api.db.Delete(&labelEntities); result.Error != nil {
		log.Printf("Delete Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context := utils.MakeResponseResultSuccess(nil)
	c.JSON(http.StatusOK, context)
}
