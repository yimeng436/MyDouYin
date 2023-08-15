package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/yimeng436/MyDouYin/pkg/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"usersvr/config"
	"usersvr/log"
	"usersvr/middleware/consul"
	"usersvr/middleware/db"
	"usersvr/service"
)

func main() {

	Init()
	defer log.Sync()
	defer db.CloseDB()
	if err := Run(); err != nil {
		log.Errorf("UserSvr run err:%v", err)
	}
}

func Init() {

	err := config.Init()
	if err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	log.InitLog()
	log.Infof("config loadding success")

}
func Run() error {
	globalConfig := config.GetGlobalConfig()
	svrconfig := globalConfig.SvrConfig
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", svrconfig.Host, svrconfig.Port))
	if err != nil {
		log.Fatalf("端口监听异常", err)
		return fmt.Errorf("端口监听异常 %v", err)
	}

	server := grpc.NewServer()
	userserver := &service.UserService{}

	pb.RegisterUserServiceServer(server, userserver)

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	register := consul.NewRegister(globalConfig.ConsulConfig.Host, globalConfig.ConsulConfig.Port)
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	err = register.Register(globalConfig.SvrConfig.Host, globalConfig.SvrConfig.Port,
		globalConfig.Name, globalConfig.ConsulConfig.Tags, serviceID)

	if err != nil {
		log.Fatal("consul.Register error: ", zap.Error(err))
		return fmt.Errorf("consul.Register error: ", zap.Error(err))
	}
	log.Info("Init Consul Register success")
	// 启动
	log.Infof("TikTokLite.user_svr listening on %s:%d", config.GetGlobalConfig().
		SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port)
	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("grpc 开启失败")
		}
	}()

	quit := make(chan os.Signal)
	//告诉操作系统，当接收到 SIGINT 或 SIGTERM 这个信号时候会将信号传给通道
	//这样做的好处是可以在操作系统接收到 中断或者关闭时候优雅的执行后续程序
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//等待接收信号，阻塞
	<-quit
	// 服务终止，注销 consul 服务
	if err = register.DeRegister(serviceID); err != nil {
		log.Info("注销失败")
		return fmt.Errorf("注销失败")
	} else {
		log.Info("注销成功")
	}
	return nil
}
