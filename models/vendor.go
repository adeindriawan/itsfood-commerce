package models

import (
	"time"
)

type Vendor struct {
	ID uint64									`gorm:"primaryKey"`
	CompanyName string				`gorm:"column:company_name" json:"company_name"`
	CompanyType string				`gorm:"column:company_type" json:"company_type"`
	Phone string							`gorm:"column:phone" json:"phone"`
	Address string						`gorm:"column:address" json:"address"`
	Village string						`gorm:"column:village" json:"village"`
	District string 					`gorm:"column:district" json:"district"`
	Regency string 						`gorm:"column:regency" json:"regency"`
	Province string 					`gorm:"column:province" json:"province"`
	PostalCode string					`gorm:"column:postal_code" json:"postal_code"`
	NPWPNumber string					`gorm:"column:npwp_number" json:"npwp_number"`
	NPWPName string						`gorm:"column:npwp_name" json:"npwp_name"`
	NPWPAddress string				`gorm:"column:npwp_address" json:"npwp_address"`
	OfficerName string 				`gorm:"column:officer_name" json:"officer_name"`
	OfficerPhone string 			`gorm:"column:officer_phone" json:"officer_phone"`
	OfficerPosition string		`gorm:"column:officer_position" json:"officer_position"`
	OfficerAddress string			`gorm:"column:officer_address" json:"officer_address"`
	OfficerIDNumber string 		`gorm:"column:officer_id_number" json:"officer_id_number"`
	PKPNumber string 					`gorm:"column:pkp_number" json:"pkp_number"`
	PKPExpiryDate time.Time		`gorm:"column:pkp_expiry_date" json:"pkp_expiry_date"`
	BankName string 					`gorm:"column:bank_name" json:"bank_name"`
	BankBranch string 				`gorm:"column:bank_branch" json:"bank_branch"`
	BankAccountNumber string	`gorm:"column:bank_account_number" json:"bank_account_number"`
	BankAccountName string 		`gorm:"column:bank_account_name" json:"bank_account_name"`
	VendorMinOrderAmount uint `gorm:"column:vendor_min_order_amount;default:0" json:"vendor_min_order_amount"`
	VendorMinOrderQty uint		`gorm:"column:vendor_min_order_qty;default:1" json:"vendor_min_order_qty"`
	VendorDeliveryCost uint 	`gorm:"column:vendor_delivery_cost;default:0" json:"vendor_delivery_cost"`
	VendorServiceCharge uint 	`gorm:"column:vendor_service_charge;default:0" json:"vendor_service_charge"`
	VendorMargin float64			`gorm:"column:vendor_margin;not null;default:10" json:"vendor_margin"`
	VendorNoteForMenus string	`gorm:"column:vendor_note_for_menus" json:"vendor_note_for_menus"`
	VendorTelegramID string 	`gorm:"column:vendor_telegram_id" json:"vendor_telegram_id"`
	Status string 						`gorm:"column:status;not null" json:"status"`
	UserID uint64 						`gorm:"column:user_id;not null" json:"user_id"`
	CreatedBy string 					`gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt time.Time 			`gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time				`gorm:"column:updated_at;autoUpdateTime:false" json:"updated_at"`			
}