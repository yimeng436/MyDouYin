package service

import (
	pb "github.com/yimeng436/MyDouYin/pkg/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (s UserService) CheckPassword(request pb.CheckPassWordRequest) pb.CheckPassWordResponse {

}
