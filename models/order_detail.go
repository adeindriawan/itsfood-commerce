package models

import (
	"time"
	"github.com/adeindriawan/itsfood-commerce/services"
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
	Note string					`gorm:"column:note" json:"note"`
	Status string 			`gorm:"column:status;not null" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time	`gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string 		`gorm:"column:created_by;not null" json:"created_by"`
}

type OrderDetailDump struct {
	ID uint64 					`gorm:"primaryKey" json:"id"`
	SourceID uint64			`gorm:"column:source_id;not null" json:"source_id"`
	OrderDetail OrderDetail `gorm:"foreignKey:SourceID" json:"order_detail"`
	OrderID uint64 			`gorm:"column:order_id;not null" json:"order_id"`
	Order Order 				`json:"order"`
	MenuID uint64 			`gorm:"column:menu_id;not null" json:"menu_id"`
	Menu Menu						`json:"menu"`
	Qty uint 						`gorm:"column:qty;not null" json:"qty"`
	Price uint64 				`gorm:"column:price;not null" json:"price"`
	COGS uint64 				`gorm:"column:cogs;not null" json:"cogs"`
	Note string					`gorm:"column:note" json:"note"`
	Status string 			`gorm:"column:status;not null" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time	`gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string 		`gorm:"column:created_by;not null" json:"created_by"`
}

func (OrderDetailDump) TableName() string {
	return "__order_details"
}

func UpdateOrderDetail(params map[string]interface{}, update map[string]interface{}) {
	var orderDetails []OrderDetail
	services.DB.Find(&orderDetails, params)

	for _, item := range orderDetails {
		orderDetailDump := OrderDetailDump{
			SourceID: item.ID,
			OrderID: item.OrderID,
			MenuID: item.MenuID,
			Qty: item.Qty,
			Price: item.Price,
			COGS: item.COGS,
			Note: item.Note,
			Status: item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			CreatedBy: item.CreatedBy,
		}
		services.DB.Create(&orderDetailDump)
	}
	services.DB.Model(&orderDetails).Updates(update)
}