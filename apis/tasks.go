package apis

import (
	"log"
	"net/http"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/database/entities"
	"github.com/pc01pc013/task-management/models/tasks"
	"github.com/pc01pc013/task-management/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TasksApi struct {
	db *gorm.DB
}

func NewTasksApi(dbInstance *gorm.DB) *TasksApi {
	return &TasksApi{db: dbInstance}
}

func (api *TasksApi) List(c *gin.Context) {
	var taskEntities []entities.Task
	api.db.Preload("Label").Find(&taskEntities)

	var result []tasks.TaskListRes
	From(taskEntities).
		SelectT(func(i entities.Task) tasks.TaskListRes {
			var labelstr []string
			From(i.Label).Select(func(i interface{}) interface{} { return i.(*entities.Label).Name }).ToSlice(&labelstr)

			return tasks.TaskListRes{
				ID:          int(i.ID),
				Title:       i.Title,
				Description: i.Description,
				StartTime:   i.StartTime,
				EndTime:     i.EndTime,
				Priority:    i.Priority,
				State:       i.State,
				Labels:      labelstr,
			}
		}).ToSlice(&result)

	context := utils.MakeResponseResultSuccess(result)

	c.JSON(http.StatusOK, context)
}

func (api *TasksApi) Create(c *gin.Context) {
	var req tasks.TaskCreateReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var tasklabels []*entities.Label

	if len(req.Labels) > 0 {
		api.db.Find(&tasklabels, req.Labels)
	}

	newtask := entities.Task{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Priority:    utils.JsonNumberPointToIntPoint(req.Priority),
		State:       utils.JsonNumberPointToIntPoint(req.State),
		Label:       tasklabels,
		UserID:      1,
	}

	if result := api.db.Create(&newtask); result.Error != nil {
		log.Printf("Create Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context := utils.MakeResponseResultSuccess(newtask.ID)
	c.JSON(http.StatusOK, context)
}

func (api *TasksApi) Update(c *gin.Context) {
	var req tasks.TaskUpdateReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var taskEntities entities.Task
	reqID, _ := req.ID.Int64()
	if result := api.db.First(&taskEntities, reqID); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	taskEntities.Title = req.Title
	taskEntities.Description = req.Description
	taskEntities.StartTime = req.StartTime
	taskEntities.EndTime = req.EndTime
	taskEntities.Priority = utils.JsonNumberPointToIntPoint(req.Priority)
	taskEntities.State = utils.JsonNumberPointToIntPoint(req.State)

	var tasklabels []*entities.Label
	isUpdateLabels := len(req.Labels) > 0

	if isUpdateLabels {
		api.db.Find(&tasklabels, req.Labels)
	}

	taskEntities.Label = tasklabels

	if err := api.db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&taskEntities).Association("Label").Replace(tasklabels)
		api.db.Omit(clause.Associations).Save(&taskEntities)

		return nil
	}); err != nil {
		log.Printf("Transaction Error: %q", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context := utils.MakeResponseResultSuccess(nil)
	c.JSON(http.StatusOK, context)
}
