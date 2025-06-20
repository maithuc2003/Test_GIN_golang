package models

import (
	"time"
)

type Book struct {
	ID        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title"`
	Stock     int       `json:"stock"`
	AuthorID  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
