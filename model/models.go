package model

import (
	"time"
)

type Topic struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     string    `json:"title" gorm:"size:256;default:" validate:"max=256"`
	Content   string    `json:"content"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:updatedAt"`
}

type Message struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	TopicID   uint      `json:"topic_id" gorm:"column:topicId"`
	Topic     Topic     `gorm:"column:topicId`
	Content   string    `json:"content"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:updatedAt"`
}
