package repository

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"usersvr/log"
	"usersvr/middleware/db"
)

func isUserNameExist(username string) (bool, error) {
	db := db.GetDB()
	var count int64
	err := db.Model(&User{}).Where("user_name=?", username).Count(&count).Error

	if err != nil {
		log.Fatal("数据库错误", err)
		return false, err
	}
	if count != 0 {
		return true, errors.New("用户名已存在")
	}

	return false, nil
}

func Register(userName, password string) (*User, error) {
	exist, err := isUserNameExist(userName)
	if err != nil {
		log.Fatal("Register", err)
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("user %s exists", userName)
	}
	db := db.GetDB()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := User{
		Name:            userName,
		Password:        string(hashPassword),
		Follow:          0,
		Follower:        0,
		TotalFav:        0,
		FavCount:        0,
		Avatar:          "https://tse1-mm.cn.bing.net/th/id/R-C.d83ded12079fa9e407e9928b8f300802?rik=Gzu6EnSylX9f1Q&riu=http%3a%2f%2fwww.webcarpenter.com%2fpictures%2fGo-gopher-programming-language.jpg&ehk=giVQvdvQiENrabreHFM8x%2fyOU70l%2fy6FOa6RS3viJ24%3d&risl=&pid=ImgRaw&r=0",
		BackgroundImage: "https://tse2-mm.cn.bing.net/th/id/OIP-C.sDoybxmH4DIpvO33-wQEPgHaEq?pid=ImgDet&rs=1",
		Signature:       "test sign",
	}

	res := db.Model(User{}).Create(&user)

	if res.Error != nil {
		return nil, res.Error
	}
	zap.L().Info("create user", zap.Any("user", user))
	//go CacheSetUser(user)
	return &user, nil
}

func GetUserInfo(selectFeild interface{}) (User, error) {
	db := db.GetDB()
	var user User
	var err error
	switch selectFeild.(type) {
	case int64:
		err = db.Where("id= ?", selectFeild).First(&user).Error
	case string:
		err = db.Where("user_name= ?", selectFeild).First(&user).Error
	default:
		err = errors.New("类型错误")
	}

	if err != nil {
		return user, err
	}
	return user, nil
}
