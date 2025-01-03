package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login        string   `json:"login"`
	Password     string   `json:"-"`
	Groups       []string `json:"groups" gorm:"type:text[]"`
	RefreshToken string   `json:"-"`
}
