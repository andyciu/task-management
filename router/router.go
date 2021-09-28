package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/apis"
	"github.com/pc01pc013/task-management/database"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	dbInstance := database.GetDBInstance()

	testapi := apis.NewTestApi(dbInstance)

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/db", testapi.DbTest)

	router.GET("/repeat", testapi.Repeat)

	return router
}
