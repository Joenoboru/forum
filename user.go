package main

import (
	"forum/model"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func register(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.Error(err)
		return
	}
	if err := validator.New().Struct(user); err != nil {
		c.Error(err)
		return
	}
	// 獲取表單資料
	name := user.Name
	account := user.Account
	password := user.Password

	hashedPassword, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加密密碼失敗"})
		return
	}

	newUser := model.User{Name: name, Account: account, Password: string(hashedPassword)}

	if err := db.Create(&newUser).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "User created successfully!", "resourceId": newUser.ID})

}

func login(c *gin.Context) {
	var user model.UserLogin
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	// Check if the user exists in database
	var dbUser model.User
	if err := db.Where("account = ?", user.Account).First(&dbUser).Error; err != nil {
		c.JSON(401, gin.H{"error": "Incorrect account or password"})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Incorrect account or password"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": dbUser.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	// Sign the token
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Error signing token"})
		return
	}

	c.JSON(200, gin.H{
		"token": "Bearer " + tokenString,
	})

}
