package models

import (
	"time"
)

type Customer struct {
	ID uint64						`gorm:"primaryKey" json:"id"`
	UserID uint64				`gorm:"column:user_id;not null" json:"user_id"`
	User User						`json:"user"`
	Type string 				`gorm:"column:type;not null" json:"type"`
	UnitID uint64 			`gorm:"column:unit_id;not null" json:"unit_id"`
	Unit Unit 					`json:"unit"`
	Status string 			`gorm:"column:status;not null" json:"status"`
	CreatedBy string		`gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt time.Time	`gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}