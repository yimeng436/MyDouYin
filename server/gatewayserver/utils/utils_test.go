package utils

import (
	"context"
	"gatewaysvr/config"
	"gatewaysvr/log"
	"github.com/yimeng436/MyDouYin/pkg/pb"
	"testing"
)

func TestVideoConnection(t *testing.T) {
	config.Init()
	log.InitLog()
	client := NewUserSvrClient("usersvr")
	t.Log(client)
	if client == nil {

		t.Errorf("NewVideoSvrClient err")
	}

	resp, err := client.CheckPassword(context.Background(), &pb.CheckPassWordRequest{
		Username: "123",
		Password: "123",
	})
	t.Log(resp, err)
}
