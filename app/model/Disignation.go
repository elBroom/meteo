package model

import "github.com/jinzhu/gorm"

type Disignation struct {
	gorm.Model

	Pin   string `gorm:"not null;unique_index"`
	Name  string `gorm:"not null"`
	Color string
	Unit  string
}

func (Disignation) TableName() string {
	return "disignation"
}
