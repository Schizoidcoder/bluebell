package routes

import (
	"bluebell/controller"
	"bluebell/logger"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.POST("/signup", controller.SignUpHandler)
	return r
}
