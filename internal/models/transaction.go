package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	FromUserID uint `gorm:"index"`
	ToUserID   uint `gorm:"index"`
	Amount     int  `gorm:"not null"`
}
