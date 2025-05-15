package models

import "time"

type Subscription struct {
	ID        uint   `gorm:"primaryKey"`
	Email     uint   `gorm:"not null"`
	City      string `gorm:"not null"`
	Frequency string `gorm:"not null"`
	Token     string `gorm:"not null"`
	Confirmed bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
