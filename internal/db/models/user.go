package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login        string         `json:"login" gorm:"unique;not null"`
	Password     string         `json:"password"`
	Groups       pq.StringArray `json:"groups" gorm:"type:text[]"`
	RefreshToken string         `json:"-"`
}
