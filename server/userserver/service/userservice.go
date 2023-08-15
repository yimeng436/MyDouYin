package service

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	pb "github.com/yimeng436/MyDouYin/pkg/pb"
	"golang.org/x/crypto/bcrypt"
	"time"
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
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(info.Password))
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
