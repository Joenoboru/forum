package model

import (
	"time"
)

type Topic struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserID    uint      `json:"user_id" gorm:"column:userId"`
	User      User      `gorm:"column:userId`
	Title     string    `json:"title" gorm:"size:256;default:" validate:"max=256"`
	Content   string    `json:"content"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:updatedAt"`
}

type Message struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserID    uint      `json:"user_id" gorm:"column:userId"`
	User      User      `gorm:"column:userId`
	TopicID   uint      `json:"topic_id" gorm:"column:topicId"`
	Topic     Topic     `gorm:"column:topicId`
	Content   string    `json:"content"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:updatedAt"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Account   string    `json:"account" gorm:"size:255;not null;unique"`
	Password  string    `json:"password" gorm:"size:100;not null"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column:updatedAt"`
}

type UserLogin struct {
	Account  string `json:"account" `
	Password string `json:"password" `
}
