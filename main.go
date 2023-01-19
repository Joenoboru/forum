package main

import (
	"fmt"

	"forum/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

const (
	UserName string = "root"
	Password string = ""
	Addr     string = "127.0.0.1"
	Port     int    = 3306
	Database string = "goDB"
)

func init() {
	// 設定資料庫連線
	var err error
	//組合sql連線字串
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", UserName, Password, Addr, Port, Database)
	db, err = gorm.Open(mysql.Open(addr))
	if err != nil {
		panic("failed to connect database")
	}
	// 自動建立資料表
	db.AutoMigrate(&model.Message{})
	db.AutoMigrate(&model.Topic{})
}

func main() {
	r := gin.Default()
	r.Use(ErrorMiddleware())
	r.POST("/topics", createTopic)
	r.PUT("/topics/:id", updateTopic)
	r.POST("/topics/:id/messages", createMessage)
	r.PUT("/messages/:id", updateMessage)
	r.Run()
}

// 錯誤處理
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": c.Errors,
			})
		}
	}
}
