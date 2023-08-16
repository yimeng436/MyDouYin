package route

import (
	"gatewaysvr/config"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {

	if config.GetGlobalConfig().SvrConfig.Mode == gin.ReleaseMode {
		// gin设置成发布模式：gin不在终端输出日志
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.Default()
	douyin := engine.Group("/douyin")

	{
		UserRoute(douyin)
	}
	return engine
}
