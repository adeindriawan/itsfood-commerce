package models

import (
	"time"
)

type OrderDetail struct {
	ID uint64 					`gorm:"primaryKey" json:"id"`
	OrderID uint64 			`gorm:"column:order_id;not null" json:"order_id"`
	Order Order 				`json:"order"`
	MenuID uint64 			`gorm:"column:menu_id;not null" json:"menu_id"`
	Menu Menu						`json:"menu"`
	Qty uint 						`gorm:"column:qty;not null" json:"qty"`
	Price uint64 				`gorm:"column:price;not null" json:"price"`
	COGS uint64 				`gorm:"column:cogs;not null" json:"cogs"`
	Status string 			`gorm:"column:status;not null" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time	`gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string 		`gorm:"column:created_by;not null" json:"created_by"`
}