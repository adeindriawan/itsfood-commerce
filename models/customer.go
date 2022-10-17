package models

import (
	"time"
)

type Customer struct {
	ID uint64						`gorm:"primaryKey"`
	UserID uint64				`gorm:"column:user_id;not null" json:"user_id"`
	Type string 				`gorm:"column:type;not null" json:"type"`
	UnitID uint64 			`gorm:"column:unit_id;not null" json:"unit_id"`
	Status string 			`gorm:"column:status;not null" json:"status"`
	CreatedBy string		`gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt time.Time	`gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}