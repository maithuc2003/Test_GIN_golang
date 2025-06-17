package models

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BookID    uint      `json:"book_id"`
	UserID    uint      `json:"user_id"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	OrderedAt time.Time `gorm:"autoCreateTime" json:"ordered_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
