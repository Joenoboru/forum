package main

import (
	"forum/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func createMessage(c *gin.Context) {
	var message model.Message

	if err := c.BindJSON(&message); err != nil {
		c.Error(err)
		return
	}

	if err := validator.New().Struct(message); err != nil {
		c.Error(err)
		return
	}

	topicID := c.Param("id")
	var topic model.Topic

	if err := db.First(&topic, topicID).Error; err != nil {
		c.Error(err)
		return
	}

	if topic.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No topic found!"})
		return
	}

	message.TopicID = topic.ID

	if err := db.Create(&message).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Message created successfully!", "resourceId": message.ID})
}

func updateMessage(c *gin.Context) {
	var message model.Message
	messageID := c.Param("id")

	if err := db.First(&message, messageID).Error; err != nil {
		c.Error(err)
		return
	}
	if message.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No message found!"})
		return
	}

	if err := c.BindJSON(&message); err != nil {
		c.Error(err)
		return
	}

	if err := validator.New().Struct(message); err != nil {
		c.Error(err)
		return
	}

	if err := db.Save(&message).Error; err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Message updated successfully!"})
}
