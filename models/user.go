package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Position     string    `gorm:"not null" json:"position"`
	PasswordHash string    `gorm:"not null" json:"password"`
	IsActive     bool      `gorm:"default:1" json:"isActive"`
	Role         string    `gorm:"null" json:"role"`
	Company      string    `gorm:"null" json:"company"`
	Country      string    `gorm:"null" json:"country"`
	LastLogin    time.Time `gorm:"null" json:"lastLogin"`
	CreatedBy    uint      `gorm:"null" json:"createdBy"`
	UpdatedBy    uint      `gorm:"null" json:"updatedBy"`
	CreatedAt    time.Time `gorm:"null" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"null" json:"updatedAt"`
}

type UserInput struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Position string `json:"position" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserSession struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
