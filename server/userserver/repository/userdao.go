package repository

import (
	"errors"
	"usersvr/middleware/db"
)

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
