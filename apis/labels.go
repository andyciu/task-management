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
	userid, err := utils.GetUserID(c, api.db)
	if err != nil {
		log.Printf("Find user Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var labelEntities []entities.Label
	api.db.Where("user_id = ?", userid).Order("id").Find(&labelEntities)

	var result []labels.LabelListRes
	num := 0
	From(labelEntities).
		Select(func(i interface{}) interface{} {
			num++
			return labels.LabelListRes{
				ID:   int(i.(entities.Label).ID),
				Num:  num,
				Name: i.(entities.Label).Name,
			}
		}).ToSlice(&result)

	if result == nil {
		result = make([]labels.LabelListRes, 0)
	}

	context := utils.MakeResponseResultSuccess(result)

	c.JSON(http.StatusOK, context)
}

func (api *LabelsApi) Create(c *gin.Context) {
	userid, err := utils.GetUserID(c, api.db)
	if err != nil {
		log.Printf("Find user Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var req labels.LabelCreateReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	if mode := c.GetString("sysmode"); mode == "nil" {
		var count int64
		if api.db.Model(&entities.Label{}).Count(&count); count > 100 {
			context := utils.MakeResponseResultFailed("")
			c.JSON(http.StatusOK, context)
			return
		}
	}

	newlabel := entities.Label{
		Name:   req.Name,
		UserID: userid,
	}
	if result := api.db.Create(&newlabel); result.Error != nil {
		log.Printf("Create Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	context := utils.MakeResponseResultSuccess(newlabel.ID)
	c.JSON(http.StatusOK, context)
}

func (api *LabelsApi) Update(c *gin.Context) {
	var req labels.LabelUpdateReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var labelEntities entities.Label
	reqID, _ := req.ID.Int64()
	if result := api.db.First(&labelEntities, reqID); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	if labelEntities.Name != req.Name {
		labelEntities.Name = req.Name
	}

	if result := api.db.Save(&labelEntities); result.Error != nil {
		log.Printf("Save Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	context := utils.MakeResponseResultSuccess(nil)
	c.JSON(http.StatusOK, context)
}

func (api *LabelsApi) Delete(c *gin.Context) {
	var req labels.LabelDeleteReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var labelEntities entities.Label
	reqID, _ := req.ID.Int64()
	if result := api.db.First(&labelEntities, reqID); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	if result := api.db.Delete(&labelEntities); result.Error != nil {
		log.Printf("Delete Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	context := utils.MakeResponseResultSuccess(nil)
	c.JSON(http.StatusOK, context)
}
