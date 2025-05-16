package models

import "time"

type Subscription struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"not null"`
	City      string    `json:"city" gorm:"not null"`
	Frequency string    `json:"frequency" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	Confirmed bool      `json:"confirmed" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
