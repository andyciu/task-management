package apis

import (
	"log"
	"net/http"
	"strconv"

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
	userid, err := utils.GetUserID(c, api.db)
	if err != nil {
		log.Printf("Find user Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	reqTitle := c.Query("title")
	reqDescription := c.Query("description")
	reqLabels := c.QueryArray("labels") //list?labels=1&labels=2

	var taskEntities []entities.Task
	tx := api.db

	if len(reqTitle) > 0 {
		tx = tx.Where("title LIKE ?", utils.QueryCondLikeString(reqTitle))
	}
	if len(reqDescription) > 0 {
		tx = tx.Where("description LIKE ?", utils.QueryCondLikeString(reqDescription))
	}
	if len(reqLabels) > 0 {
		var labelEntities []entities.Label
		From(reqLabels).Select(func(i interface{}) interface{} {
			k, _ := strconv.Atoi(i.(string))
			return entities.Label{
				Model: gorm.Model{
					ID: uint(k),
				},
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

	if result := tx.Where("user_id = ?", userid).Preload("Label").Order("id").Find(&taskEntities); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var result []tasks.TaskListRes
	num := 0
	From(taskEntities).
		SelectT(func(i entities.Task) tasks.TaskListRes {
			var labelsnum []int

			From(i.Label).Select(func(i interface{}) interface{} { return int(i.(*entities.Label).ID) }).ToSlice(&labelsnum)
			num++
			return tasks.TaskListRes{
				ID:          int(i.ID),
				Num:         num,
				Title:       i.Title,
				Description: i.Description,
				StartTime:   i.StartTime,
				EndTime:     i.EndTime,
				Priority:    i.Priority,
				State:       i.State,
				Labels:      labelsnum,
			}
		}).ToSlice(&result)

	if result == nil {
		result = make([]tasks.TaskListRes, 0)
	}

	context := utils.MakeResponseResultSuccess(result)

	c.JSON(http.StatusOK, context)
}

func (api *TasksApi) Create(c *gin.Context) {
	userid, err := utils.GetUserID(c, api.db)
	if err != nil {
		log.Printf("Find user Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var req tasks.TaskCreateReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	if mode := c.GetString("sysmode"); mode == "nil" {
		var count int64
		if api.db.Model(&entities.Task{}).Count(&count); count > 100 {
			context := utils.MakeResponseResultFailed("")
			c.JSON(http.StatusOK, context)
			return
		}
	}

	var tasklabels []*entities.Label

	if len(req.Labels) > 0 {
		api.db.Find(&tasklabels, req.Labels)
	}

	newtask := entities.Task{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   utils.DateTimePrase(req.StartTime),
		EndTime:     utils.DateTimePrase(req.EndTime),
		Priority:    utils.JsonNumberPointToIntPoint(req.Priority),
		State:       utils.JsonNumberPointToIntPoint(req.State),
		Label:       tasklabels,
		UserID:      userid,
	}

	if result := api.db.Create(&newtask); result.Error != nil {
		log.Printf("Create Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	context := utils.MakeResponseResultSuccess(newtask.ID)
	c.JSON(http.StatusOK, context)
}

func (api *TasksApi) Update(c *gin.Context) {
	var req tasks.TaskUpdateReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var taskEntities entities.Task
	reqID, _ := req.ID.Int64()
	if result := api.db.First(&taskEntities, reqID); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
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
		tx.Model(&taskEntities).Association("Label").Clear()
		if isUpdateLabels {
			tx.Model(&taskEntities).Association("Label").Append(tasklabels)
		}
		api.db.Omit(clause.Associations).Save(&taskEntities)

		return nil
	}); err != nil {
		log.Printf("Transaction Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	context := utils.MakeResponseResultSuccess(nil)
	c.JSON(http.StatusOK, context)
}

func (api *TasksApi) Delete(c *gin.Context) {
	var req tasks.TaskDeleteReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	var taskEntities entities.Task
	reqID, _ := req.ID.Int64()
	if result := api.db.First(&taskEntities, reqID); result.Error != nil {
		log.Printf("Find Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	if result := api.db.Select(clause.Associations).Delete(&taskEntities); result.Error != nil {
		log.Printf("Delete Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("")
		c.JSON(http.StatusOK, context)
		return
	}

	context := utils.MakeResponseResultSuccess(nil)
	c.JSON(http.StatusOK, context)
}
