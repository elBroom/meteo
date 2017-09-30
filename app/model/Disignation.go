package model

type Disignation struct {
	ID    uint   `gorm:"primary_key"`
	Pin   string `gorm:"not null;unique_index"`
	Name  string `gorm:"not null"`
	Color string
	Unit  string
}

func (Disignation) TableName() string {
	return "meteo_disignation"
}
