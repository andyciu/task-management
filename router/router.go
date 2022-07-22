package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/apis"
	"github.com/pc01pc013/task-management/auth"
	"github.com/pc01pc013/task-management/database"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.New(auth.CorsConfig()))
	router.Use(auth.SysMode())

	dbInstance := database.GetDBInstance()

	labelsApi := apis.NewLabelsApi(dbInstance)
	tasksApi := apis.NewTasksApi(dbInstance)
	authApi := apis.NewAuthApi(dbInstance)
	userApi := apis.NewUserApi(dbInstance)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authApi.Login)
	}

	apisGroup := router.Group("/apis")
	apisGroup.Use(auth.JWTAuthMiddleware())
	{
		labelsRoute := apisGroup.Group("/labels")
		{
			labelsRoute.GET("/list", labelsApi.List)
			labelsRoute.POST("/create", labelsApi.Create)
			labelsRoute.POST("/update", labelsApi.Update)
			labelsRoute.POST("/deleteL", labelsApi.Delete)
		}
		tasksRoute := apisGroup.Group("/tasks")
		{
			tasksRoute.GET("/list", tasksApi.List)
			tasksRoute.POST("/create", tasksApi.Create)
			tasksRoute.POST("/update", tasksApi.Update)
			tasksRoute.POST("/deleteL", tasksApi.Delete)
		}
		userRoute := apisGroup.Group("/user")
		{
			userRoute.GET("getNickName", userApi.GetNickName)
		}
	}

	return router
}
