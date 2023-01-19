package main

import (
	"forum/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func createTopic(c *gin.Context) {
	var topic model.Topic

	if err := c.BindJSON(&topic); err != nil {
		c.Error(err)
		return
	}
	if err := validator.New().Struct(topic); err != nil {
		c.Error(err)
		return
	}
	if err := db.Create(&topic).Error; err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Topic created successfully!", "resourceId": topic.ID})
}

func updateTopic(c *gin.Context) {
	var topic model.Topic
	topicID := c.Param("id")

	if err := db.First(&topic, topicID).Error; err != nil {
		c.Error(err)
		return
	}

	if topic.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No topic found!"})
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
