package service

import (
	"context"
	"github.com/yimeng436/MyDouYin/pkg/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (s UserService) CheckPassword(ctx context.Context, request *pb.CheckPassWordRequest) (*pb.CheckPassWordResponse, error) {

}
