package main

import (
	"fmt"
	"forum/model"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	// 設定資料庫連線
	var err error

	enverr := godotenv.Load()
	if enverr != nil {
		panic("Error loading .env file")
	}
	UserName := os.Getenv("DB_USER")
	Password := os.Getenv("DB_PASSWORD")
	Host := os.Getenv("DB_HOST")
	Database := os.Getenv("DB_NAME")
	Port := os.Getenv("DB_PORT")
	//組合sql連線字串
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", UserName, Password, Host, Port, Database)
	db, err = gorm.Open(mysql.Open(addr), &gorm.Config{})
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
