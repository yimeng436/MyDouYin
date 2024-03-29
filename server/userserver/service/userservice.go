package service

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	pb "github.com/yimeng436/MyDouYin/pkg/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"usersvr/log"
	"usersvr/repository"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

type JWTClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

func (UserService) CheckPassword(ctx context.Context, request *pb.CheckPassWordRequest) (*pb.CheckPassWordResponse, error) {
	username := request.Username
	password := request.Password

	//这里 持久化层做了一个函数复用，因为我们经常可能要通过某个字段查找用户
	//所以可以将这个通用函数的类型定义为interface，用类型去调用
	//当然这种做法只适合所有类型都不相同的查找
	info, err := repository.GetUserInfo(username)
	if err != nil {
		return nil, err
	}

	//必须第一个参数是数据库加密后的数据
	err = bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	token, err := GenToken(info.Id, username)
	if err != nil {
		return nil, err
	}
	response := &pb.CheckPassWordResponse{
		UserId: info.Id,
		Token:  token,
	}
	return response, nil
}

func (UserService) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userName := request.Username
	password := request.Password

	user, err := repository.Register(userName, password)

	if err != nil {
		return nil, err
	}

	token, err := GenToken(user.Id, userName)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		UserId: user.Id,
		Token:  token,
	}, nil
}

func (UserService) GetUserInfoList(context context.Context, request *pb.GetUserInfoListRequest) (*pb.GetUserInfoListResponse, error) {
	info, err := repository.GetUserInfo(request.UserId)
	if err != nil {
		log.Fatal("repository.GetUserInfo", err)
		return nil, err
	}

	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoList not implemented")
}
func (UserService) GetUserInfo(context context.Context, request *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UserService) GetUserInfoDict(context context.Context, request *pb.GetUserInfoDictRequest) (*pb.GetUserInfoDictResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoDict not implemented")
}

func GenToken(id int64, username string) (string, error) {
	claims := &JWTClaims{
		UserId:   id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "server",
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * 30)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte("douyin"))
	if err != nil {
		return "", err
	}

	return signedString, nil
}
