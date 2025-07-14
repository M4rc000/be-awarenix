// models/campaign.go
package models

import (
	"time"
)

type Campaign struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string     `gorm:"type:varchar(50);not null" json:"name" binding:"required"`
	LaunchDate       time.Time  `gorm:"type:datetime;not null" json:"launch_date" binding:"required"` // Pastikan ini time.Time
	SendEmailBy      *time.Time `gorm:"type:datetime;null" json:"send_email_by,omitempty"`
	GroupID          int64      `gorm:"not null;index" json:"group_id" binding:"required"`
	EmailTemplateID  int64      `gorm:"not null;index" json:"email_template_id" binding:"required"`
	LandingPageID    int64      `gorm:"not null;index" json:"landing_page_id" binding:"required"`
	SendingProfileID int64      `gorm:"not null;index" json:"sending_profile_id" binding:"required"`
	URL              string     `gorm:"type:varchar(255);not null" json:"url" binding:"required,url"`
	CreatedBy        int64      `gorm:"type:tinyint(3);null" json:"created_by"`
	CreatedAt        time.Time  `gorm:"type:datetime;null" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:datetime;null" json:"updated_at"`
	Status           string     `gorm:"type:varchar(30);null" json:"status"`
	CompletedDate    *time.Time `gorm:"type:datetime;null" json:"completed_date,omitempty"`
}

type CampaignRequest struct {
	Name             string  `json:"name" binding:"required"`
	LaunchDate       string  `json:"launch_date" binding:"required"` // String untuk JSON dari frontend
	SendEmailBy      *string `json:"send_email_by,omitempty"`
	GroupID          int64   `json:"group_id" binding:"required"`
	EmailTemplateID  int64   `json:"email_template_id" binding:"required"`
	LandingPageID    int64   `json:"landing_page_id" binding:"required"`
	SendingProfileID int64   `json:"sending_profile_id" binding:"required"`
	URL              string  `json:"url" binding:"required,url"`
	CreatedBy        int64   `json:"created_by"`
}

type CampaignResponse struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	LaunchDate       time.Time  `json:"launch_date"`
	SendEmailBy      *time.Time `json:"send_email_by,omitempty"`
	GroupID          int64      `json:"group_id"`
	EmailTemplateID  int64      `json:"email_template_id"`
	LandingPageID    int64      `json:"landing_page_id"`
	SendingProfileID int64      `json:"sending_profile_id"`
	URL              string     `json:"url"`
	CreatedBy        int64      `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Status           string     `json:"status"`
	CompletedDate    *time.Time `json:"completed_date,omitempty"`
}

type NewCampaignResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}
