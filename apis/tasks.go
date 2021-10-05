package apis

import (
	"encoding/json"
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
	var req tasks.TaskListReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var taskEntities []entities.Task
	tx := api.db

	if req.Title != nil {
		tx = tx.Where("title LIKE ?", utils.QueryCondLikeString(*req.Title))
	}
	if req.Description != nil {
		tx = tx.Where("description LIKE ?", utils.QueryCondLikeString(*req.Description))
	}
	if len(req.Labels) > 0 {
		var labelEntities []entities.Label
		From(req.Labels).Select(func(i interface{}) interface{} {
			k, _ := i.(json.Number).Int64()
			return entities.Label{
				ID: uint(k),
			}
		}).ToSlice(&labelEntities)

		//select tasks.id
		type APITask struct {
			ID uint
		}
		var tasknumApi []APITask
		api.db.Model(labelEntities).Association("Task").Find(&tasknumApi)

		var tasknumInt []uint
		From(tasknumApi).Select(func(i interface{}) interface{} { return i.(APITask).ID }).ToSlice(&tasknumInt)

		tx = tx.Where("id IN (?)", tasknumInt)
	}

	if result := tx.Preload("Label").Find(&taskEntities); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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

func (api *TasksApi) Delete(c *gin.Context) {
	var req tasks.TaskDeleteReq

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

	if result := api.db.Select(clause.Associations).Delete(&taskEntities); result.Error != nil {
		log.Printf("Delete Error: %q", result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context := utils.MakeResponseResultSuccess(nil)
	c.JSON(http.StatusOK, context)
}
