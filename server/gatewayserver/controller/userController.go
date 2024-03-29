package controller

import (
	"gatewaysvr/config"
	"gatewaysvr/log"
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/gin-gonic/gin"
	"github.com/yimeng436/MyDouYin/pkg/pb"
	"go.uber.org/zap"
)

func Login(context *gin.Context) {

	request := new(pb.CheckPassWordRequest)
	err := context.ShouldBind(&request)
	if err != nil {
		log.Fatal("请求参数错误")
		return
	}

	client := utils.GetUserServiceClient()
	rep, err := client.CheckPassword(context, request)
	if err != nil {
		zap.L().Error("login error", zap.Error(err))
		response.Fail(context, err.Error(), nil)
		return
	}
	response.Success(context, "success", rep)
}

func Register(context *gin.Context) {
	var request pb.RegisterRequest
	err := context.ShouldBind(&request)

	if err != nil {
		log.Fatal("请求参数错误")
		return
	}

	resp, err := utils.NewUserSvrClient(config.GetGlobalConfig().UserSvrName).Register(context, &request)
	if err != nil {
		log.Error("register error", err)
		response.Fail(context, err.Error(), nil)
		return
	}
	response.Success(context, "success", resp)
}
