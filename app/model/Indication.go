package model

import "github.com/jinzhu/gorm"

type Indication struct {
	gorm.Model

	Pin   string  `gorm:"not null;index"`
	Value float32 `gorm:"not null"`
}

func (Indication) TableName() string {
	return "indication"
}
