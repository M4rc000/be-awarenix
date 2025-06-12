package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Position     string    `gorm:"not null" json:"position"`
	PasswordHash string    `gorm:"not null" json:"password"`
	IsActive     bool      `gorm:"default:1" json:"isActive"`
	CreatedBy    uint      `json:"createdBy"`
	UpdatedBy    uint      `json:"updatedBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserSession struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
