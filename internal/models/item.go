package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Price int    `gorm:"not null"`
}
