package models

import (
	"time"
	"github.com/adeindriawan/itsfood-commerce/services"
)

type Order struct {
	ID uint64 						`gorm:"primaryKey" json:"id"`
	OrderedBy uint64 			`gorm:"column:ordered_by;not null" json:"ordered_by"`
	Customer Customer			`gorm:"foreignKey:OrderedBy" json:"customer"`
	OrderedFor time.Time	`gorm:"column:ordered_for;not null" json:"ordered_for"`
	OrderedTo string 			`gorm:"column:ordered_to;not null" json:"ordered_to"`
	NumOfMenus uint 			`gorm:"column:num_of_menus;not null" json:"num_of_menus"`
	QtyOfMenus uint 			`gorm:"column:qty_of_menus;not null" json:"qty_of_menus"`
	Amount uint64 				`gorm:"column:amount;not null" json:"amount"`
	OrderDetail []OrderDetail `json:"order_detail"`
	Purpose string 				`gorm:"column:purpose;not null" json:"purpose"`
	Activity string 			`gorm:"column:activity;not null" json:"activity"`
	SourceOfFund string		`gorm:"column:source_of_fund;not null" json:"source_of_fund"`
	PaymentOption string	`gorm:"column:payment_option;not null" json:"payment_option"`
	Info string	 					`gorm:"info;not null" json:"info"`
	Status string 				`gorm:"column:status;not null" json:"status"`
	CreatedAt time.Time		`gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time		`gorm:"column:updated_at;not null" json:"updated_at"`
	CreatedBy string 			`gorm:"created_by;not null" json:"created_by"`
}

type OrderDump struct {
	ID uint64 						`gorm:"primaryKey" json:"id"`
	SourceID uint64 			`gorm:"column:source_id;not null" json:"source_id"`
	Order Order 					`gorm:"foreignKey:SourceID" json:"order"`
	OrderedBy uint64 			`gorm:"column:ordered_by;not null" json:"ordered_by"`
	Customer Customer			`gorm:"foreignKey:OrderedBy" json:"customer"`
	OrderedFor time.Time	`gorm:"column:ordered_for;not null" json:"ordered_for"`
	OrderedTo string 			`gorm:"column:ordered_to;not null" json:"ordered_to"`
	NumOfMenus uint 			`gorm:"column:num_of_menus;not null" json:"num_of_menus"`
	QtyOfMenus uint 			`gorm:"column:qty_of_menus;not null" json:"qty_of_menus"`
	Amount uint64 				`gorm:"column:amount;not null" json:"amount"`
	Purpose string 				`gorm:"column:purpose;not null" json:"purpose"`
	Activity string 			`gorm:"column:activity;not null" json:"activity"`
	SourceOfFund string		`gorm:"column:source_of_fund;not null" json:"source_of_fund"`
	PaymentOption string	`gorm:"column:payment_option;not null" json:"payment_option"`
	Info string	 					`gorm:"info;not null" json:"info"`
	Status string 				`gorm:"column:status;not null" json:"status"`
	CreatedAt time.Time		`gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time		`gorm:"column:updated_at;autoUpdateTime:false" json:"updated_at"`
	CreatedBy string 			`gorm:"created_by;not null" json:"created_by"`
}

func (OrderDump) TableName() string {
	return "__orders"
}

func UpdateOrder(params map[string]interface{}, update map[string]interface{}) {
	var orders []Order
	services.DB.Find(&orders, params)
	// create dump
	for _, item := range orders {
		orderDump := OrderDump{
			SourceID: item.ID,
			OrderedBy: item.OrderedBy,
			OrderedFor: item.OrderedFor,
			OrderedTo: item.OrderedTo,
			NumOfMenus: item.NumOfMenus,
			QtyOfMenus: item.QtyOfMenus,
			Amount: item.Amount,
			Purpose: item.Purpose,
			Activity: item.Activity,
			SourceOfFund: item.SourceOfFund,
			PaymentOption: item.PaymentOption,
			Info: item.Info,
			Status: item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			CreatedBy: item.CreatedBy,
		}
		services.DB.Create(&orderDump)
	}
	// update record
	services.DB.Model(&orders).Updates(update)
}