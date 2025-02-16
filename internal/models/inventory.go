package models

import "gorm.io/gorm"

// Купленные товары
type Inventory struct {
	gorm.Model
	UserID   uint   `gorm:"index"`
	ItemName string `gorm:"not null"`
	Quantity int    `gorm:"default:1"`
}
