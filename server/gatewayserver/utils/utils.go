package utils

import (
	"context"
	"fmt"
	"gatewaysvr/config"
	"gatewaysvr/log"
	"github.com/yimeng436/MyDouYin/pkg/pb"
	"google.golang.org/grpc"
	// 必须要导入这个包，否则grpc会报错
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"time"
)

var (
	UserSvrClient pb.UserServiceClient
)

func InitSvrConn() {
	UserSvrClient = NewUserSvrClient(config.GetGlobalConfig().SvrConfig.UserSvrName)
}

func NewUserSvrClient(svrName string) pb.UserServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewUserServiceClient(conn)
}

func NewSvrConn(svrName string) (*grpc.ClientConn, error) {
	consulInfo := config.GetGlobalConfig().ConsulConfig
	fmt.Println(consulInfo)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(
		ctx,
		//这里必须要引入上面github.com/mbobakov/grpc-consul-resolver 这个包
		//因为底层是通过这个解析器将这个请求转为http请求找到对应的服务
		//只有把这个包引进来了，才会执行这个包里面的consul.go的init函数，将解析器注入
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, svrName),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Errorf("NewSvrConn with svrname %s err:%v", svrName, err)
		return nil, err
	}
	log.Info("NewSvrConn success")
	return conn, nil
}
func GetUserServiceClient() pb.UserServiceClient {
	return UserSvrClient
}
