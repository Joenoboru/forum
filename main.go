package main

import (
	"flag"
	"fmt"
	"forum/model"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	envSetting := flag.String("env", ".env", "environment setting")
	flag.Parse()
	var envFile string
	if *envSetting == "docker" {
		envFile = ".env.docker"
	} else {
		envFile = ".env"
	}
	enverr := godotenv.Load(envFile)
	if enverr != nil {
		panic("Error loading .env file")
	}
	UserName := os.Getenv("DB_USER")
	Password := os.Getenv("DB_PASSWORD")
	Host := os.Getenv("DB_HOST")
	Database := os.Getenv("DB_NAME")
	Port := os.Getenv("DB_PORT")
	//組合sql連線字串
	var err error
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", UserName, Password, Host, Port, Database)
	db, err = gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 自動建立資料表
	db.AutoMigrate(&model.Message{})
	db.AutoMigrate(&model.Topic{})
	db.AutoMigrate(&model.User{})
}

func main() {
	r := gin.Default()
	r.Use(ErrorMiddleware())
	r.Use(cors.Default()) //使用 POST 方法時會有預先驗證的請求，這個請求被稱為 CORS（跨來源資源共用）預先驗證請求。可以配置服務器以回應預先驗證請求，並使用適當的 Access-Control-Allow-Origin 頭。此頭告訴瀏覽器哪些源是允許向 API 發送請求的。

	// 為特定路由群組設定中間件檢查token
	auth := r.Group("", AuthMiddleware())
	{
		auth.POST("/topics", createTopic)
		auth.PUT("/topics/:id", updateTopic)
		auth.POST("/topics/:id/messages", createMessage)
		auth.PUT("/messages/:id", updateMessage)
	}
	r.POST("/register", register)
	r.POST("/login", login)
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

// 檢查token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is empty"})
			c.Abort()
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		token := parts[1]
		tokenParseResult, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}
		if claims, ok := tokenParseResult.Claims.(jwt.MapClaims); ok && tokenParseResult.Valid {
			c.Set("userId", claims["userId"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		c.Next()
	}
}
