package models

import (
	"time"
)

type Cost struct {
	ID uint64 							`gorm:"primaryKey" json:"id"`
	OrderDetailID uint64 		`gorm:"column:order_detail_id;not null" json:"order_detail_id"`
	OrderDetail OrderDetail `json:"order_detail"`
	Amount uint							`gorm:"column:amount;not null" json:"amount"`
	Reason string 					`gorm:"column:reason;not null" json:"reason"`
	Status string 					`gorm:"column:status;not null" json:"status"`
	CreatedAt time.Time 		`gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time			`gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string 				`gorm:"column:created_by;not null" json:"created_by"`
}

type CostDump struct {
	ID uint64 							`gorm:"primaryKey" json:"id"`
	SourceID uint64 				`gorm:"column:source_id;not null" json:"source_id"`
	OrderDetailID uint64 		`gorm:"column:order_detail_id;not null" json:"order_detail_id"`
	OrderDetail OrderDetail `json:"order_detail"`
	Amount uint							`gorm:"column:amount;not null" json:"amount"`
	Reason string 					`gorm:"column:reason;not null" json:"reason"`
	Status string 					`gorm:"column:status;not null" json:"status"`
	CreatedAt time.Time 		`gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time			`gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string 				`gorm:"column:created_by;not null" json:"created_by"`
}

func (CostDump) TableName() string {
	return "__costs"
}