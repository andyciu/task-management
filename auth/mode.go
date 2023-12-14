package auth

import (
	"os"

	"github.com/gin-gonic/gin"
)

func SysMode() gin.HandlerFunc {
	return func(context *gin.Context) {
		mode := os.Getenv("MODE")
		if mode == "" {
			context.Set("sysmode", "nil")
		}
	}
}
