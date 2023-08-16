package db

import (
	"fmt"
	"gorm.io/gorm/logger"
	"usersvr/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	db     *gorm.DB
	dbonce sync.Once
)

func OpenDb() {
	dbConfig := config.GetGlobalConfig().DbConfig

	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.Username,
		dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	var er error
	db, er = gorm.Open(mysql.Open(connArgs), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if er != nil {
		fmt.Println("数据库连接失败")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("fetch db connection err:" + err.Error())
	}
	// _ = db.AutoMigrate(&repository.User{})
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)                                        // 设置最大空闲连接
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)                                        // 设置最大打开的连接
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxIdleTime * int64(time.Second))) // 设置空闲时间为(s)
	return
}

func GetDB() *gorm.DB {
	dbonce.Do(OpenDb)
	return db
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			panic("fetch db connection err:" + err.Error())
		}
		sqlDB.Close()
	}
}
