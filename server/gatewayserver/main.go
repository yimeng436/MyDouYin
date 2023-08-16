package main

import (
	"fmt"
	"gatewaysvr/config"
	"gatewaysvr/log"
	"gatewaysvr/route"
	"gatewaysvr/utils"
	"go.uber.org/zap"
)

func main() {
	Init()
	defer log.Sync()
	engine := route.InitRoute()

	err := engine.Run(fmt.Sprintf(":%d", config.GetGlobalConfig().SvrConfig.Port))
	if err != nil {
		zap.L().Panic("Router.Run error: ", zap.Error(err))
	}

	zap.L().Sugar().Infof("listen on %s:%d", config.GetGlobalConfig().SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port)

}

func Init() {
	if err := config.Init(); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	log.InitLog()
	log.Info("log init success...")
	utils.InitSvrConn()
	log.Info("Init Grpc Server conn success")
}
