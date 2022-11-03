package models

import (
	"time"
)

type Unit struct {
	ID uint64						`gorm:"primaryKey" json:"id"`
	Name string 				`gorm:"column:name;not null" json:"name"`
	GroupID uint64			`gorm:"column:group_id;not null" json:"group_id"`
	Status string 			`gorm:"column:status;not null" json:"status"`
	CreatedAt time.Time	`gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time	`gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string		`gorm:"column:created_by;not null" json:"created_by"`
}