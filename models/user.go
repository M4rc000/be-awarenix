package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(50);not null" json:"name"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Position     string    `gorm:"type:varchar(50);not null" json:"position"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"password"`
	IsActive     int       `gorm:"type:tinyint(1);default:1" json:"isActive"`
	Role         int       `gorm:"type:tinyint(3);default:3" json:"role"`
	Company      string    `gorm:"type:varchar(50);null" json:"company"`
	Country      string    `gorm:"type:varchar(50);null" json:"country"`
	LastLogin    time.Time `gorm:"type:datetime;null" json:"lastLogin"`
	CreatedAt    time.Time `gorm:"type:datetime;null" json:"createdAt"`
	CreatedBy    int       `gorm:"type:tinyint(3);null" json:"createdBy"`
	UpdatedAt    time.Time `gorm:"type:datetime;null" json:"updatedAt"`
	UpdatedBy    int       `gorm:"type:tinyint(3);null" json:"updatedBy"`
}

type FullUserLoginData struct {
	User
	RoleName string `json:"role_name"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Status   string `json:"status"`
}
type UserLoginResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Position  string    `json:"position"`
	Role      int       `json:"role"`
	RoleName  string    `json:"role_name"`
	Company   string    `json:"company"`
	Country   string    `json:"country"`
	IsActive  int       `json:"isActive"`
	LastLogin time.Time `json:"last_login"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserTable struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Position      string    `json:"position"`
	IsActive      bool      `json:"isActive"`
	Role          int       `json:"role"`
	RoleName      string    `json:"roleName"`
	Company       string    `json:"company"`
	Country       string    `json:"country"`
	LastLogin     time.Time `json:"lastLogin"`
	CreatedBy     int       `json:"createdBy"`
	UpdatedBy     int       `json:"updatedBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedByName string    `json:"createdByName"`
	UpdatedByName string    `json:"updatedByName"`
}

type CreateUserInput struct {
	Name      string    `json:"name"     binding:"required"`
	Email     string    `json:"email"    binding:"required,email"`
	Position  string    `json:"position" binding:"required"`
	Company   string    `json:"company"`
	Role      int       `json:"role"`
	Password  string    `json:"password" binding:"required,min=6"`
	CreatedAt time.Time `gorm:"null" json:"createdAt"`
	CreatedBy int       `gorm:"null" json:"createdBy"`
}

type UpdateUserInput struct {
	Name      string    `json:"name"     binding:"required"`
	Email     string    `json:"email"    binding:"required,email"`
	Position  string    `json:"position" binding:"required"`
	Role      int       `json:"role"`
	Company   string    `json:"company"`
	IsActive  int       `json:"isActive"`
	Password  string    `json:"password"`
	UpdatedAt time.Time `gorm:"null"`
	UpdatedBy int       `gorm:"null" json:"updatedBy"`
}

type UserSession struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetUserSession struct {
	ID uint `json:"user_id" gorm:"primaryKey"`
}
