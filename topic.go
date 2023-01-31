package main

import (
	"forum/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func createTopic(c *gin.Context) {
	var topic model.Topic
	userId := uint(c.MustGet("userId").(float64))

	if err := c.BindJSON(&topic); err != nil {
		c.Error(err)
		return
	}

	if err := validator.New().Struct(topic); err != nil {
		c.Error(err)
		return
	}

	topic.UserID = userId
	if err := db.Create(&topic).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Topic created successfully!", "resourceId": topic.ID})
}

func updateTopic(c *gin.Context) {
	var topic model.Topic
	topicID := c.Param("id")
	userId := uint(c.MustGet("userId").(float64))

	if err := db.First(&topic, topicID).Error; err != nil {
		c.Error(err)
		return
	}

	if topic.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No topic found!"})
		return
	}

	if topic.UserID != userId {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusUnauthorized, "message": "Wrong user!"})
		return
	}

	if err := c.BindJSON(&topic); err != nil {
		c.Error(err)
		return
	}

	if err := validator.New().Struct(topic); err != nil {
		c.Error(err)
		return
	}

	if err := db.Save(&topic).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Topic updated successfully!"})
}
