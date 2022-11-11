package models

import (
	"time"
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