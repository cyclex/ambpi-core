package model

import "github.com/jinzhu/gorm"

type Program struct {
	gorm.Model
	Retail    string `gorm:"retail" json:"retail"`
	Status    bool   `gorm:"status" json:"status"`
	StartDate int64  `gorm:"start_date" json:"startDate"`
	EndDate   int64  `gorm:"end_date" json:"endDate"`
}
