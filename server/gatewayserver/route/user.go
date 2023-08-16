package route

import (
	"gatewaysvr/controller"
	"github.com/gin-gonic/gin"
)

func UserRoute(engine *gin.RouterGroup) {
	userRoute := engine.Group("/user")
	userRoute.POST("/login", controller.Login)
}
