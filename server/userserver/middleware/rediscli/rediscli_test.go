package rediscli

import (
	"context"
	"fmt"
	"testing"
	"time"
	"usersvr/config"
	"usersvr/log"
)

func TestCloseRedis(t *testing.T) {
	if err := config.Init(); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	log.InitLog()
	log.Info("log init success...")
	cli := GetRedisCli()
	result, err := cli.Set(context.Background(), "go_test", "asd", time.Second*30).Result()
	if err != nil {
		log.Fatalf("redis 操作异常")
	}
	s, err := cli.Ping(context.Background()).Result()
	fmt.Printf(result)
	fmt.Println(s)
}
