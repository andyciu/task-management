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
