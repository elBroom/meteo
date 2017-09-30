package model

import "time"

type Indication struct {
	ID uint `gorm:"primary_key"`

	Pin        string  `gorm:"not null;index"`
	Value      float32 `gorm:"not null"`
	CreateDate time.Time
}

func (Indication) TableName() string {
	return "meteo_indication"
}
