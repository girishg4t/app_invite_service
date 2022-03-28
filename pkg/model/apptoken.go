package model

import "time"

// AppToken to store the token status
type AppToken struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Token    string    `json:"token"`
	ExpDate  time.Time `json:"exp_date"`
	IsActive bool      `json:"is_active" gorm:"default:true"`
}
