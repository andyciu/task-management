package apis

import (
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
				StartTime:   &i.StartTime,
				EndTime:     &i.EndTime,
				Priority:    i.Priority,
				State:       i.State,
				Labels:      labelstr,
			}
		}).ToSlice(&result)

	context := utils.MakeResponseResultSuccess(result)

	c.JSON(http.StatusOK, context)
}
