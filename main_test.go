package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"forum/model"

	"github.com/gin-gonic/gin"
)

func TestCreateTopic(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/topics", createTopic)

	// create request
	req, _ := http.NewRequest("POST", "/topics", bytes.NewBufferString(`{"Title":"Test Topic", "Content":"Test Content"}`))
	req.Header.Set("Content-Type", "application/json")

	// create response recorder
	w := httptest.NewRecorder()

	// execute request
	r.ServeHTTP(w, req)

	// check response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %v but got %v", http.StatusCreated, w.Code)
	}

	// teardown
	defer db.Exec("DELETE FROM topics")
}

func TestUpdateTopic(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/topics/:id", updateTopic)

	// create a topic
	var topic model.Topic
	topic.Title = "Test Topic"
	topic.Content = "Test Content"
	db.Create(&topic)

	// create request
	req, _ := http.NewRequest("PUT", "/topics/"+strconv.Itoa(int(topic.ID)), bytes.NewBufferString(`{"Title":"Updated Topic", "Content":"Updated Content"}`))
	req.Header.Set("Content-Type", "application/json")

	// create response recorder
	w := httptest.NewRecorder()

	// execute request
	r.ServeHTTP(w, req)

	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected %v but got %v", http.StatusOK, w.Code)
	}

	// check update
	var updatedTopic model.Topic
	db.First(&updatedTopic, topic.ID)
	if updatedTopic.Title != "Updated Topic" {
		t.Errorf("Expected %s but got %s", "Updated Topic", updatedTopic.Title)
	}

	// teardown
	defer db.Exec("DELETE FROM topics")
}

func TestCreateMessage(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/topics/:topicId/messages", createMessage)

	// create a topic
	var topic model.Topic
	topic.Title = "Test Topic"
	topic.Content = "Test Content"
	db.Create(&topic)
	// create request
	req, _ := http.NewRequest("POST", "/topics/"+strconv.Itoa(int(topic.ID))+"/messages", bytes.NewBufferString(`{"Content":"Test Message"}`))
	req.Header.Set("Content-Type", "application/json")

	// create response recorder
	w := httptest.NewRecorder()

	// execute request
	r.ServeHTTP(w, req)

	// check response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %v but got %v", http.StatusCreated, w.Code)
	}

	// check message
	var message model.Message
	db.Where("topicId = ?", topic.ID).First(&message)
	if message.Content != "Test Message" {
		t.Errorf("Expected %s but got %s", "Test Message", message.Content)
	}

	// teardown
	defer db.Exec("DELETE FROM topics")
	defer db.Exec("DELETE FROM messages")
}

func TestUpdateMessage(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/messages/:id", updateMessage)
	// create a topic
	var topic model.Topic
	topic.Title = "Test Topic"
	topic.Content = "Test Content"
	db.Create(&topic)

	// create a message
	var message model.Message
	message.TopicID = topic.ID
	message.Content = "Test Message"
	db.Create(&message)
	// create request
	req, _ := http.NewRequest("PUT", "/messages/"+strconv.Itoa(int(message.ID)), bytes.NewBufferString(`{"Content":"Updated Message"}`))
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()

	// execute request
	r.ServeHTTP(w, req)

	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected %v but got %v", http.StatusOK, w.Code)
	}

	// check update
	var updatedMessage model.Message
	db.First(&updatedMessage, message.ID)
	if updatedMessage.Content != "Updated Message" {
		t.Errorf("Expected %s but got %s", "Updated Message", updatedMessage.Content)
	}

	// teardown
	defer db.Exec("DELETE FROM topics")
	defer db.Exec("DELETE FROM messages")
}
