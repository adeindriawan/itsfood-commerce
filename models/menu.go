package models

import (
	"database/sql/driver"
	"time"
)

type MenuCategory string

const (
	Food MenuCategory = "Food"
	Beverage MenuCategory = "Beverage"
	Snack MenuCategory = "Snack"
	Fruit MenuCategory = "Category"
	Grocery MenuCategory = "Grocery"
	Others MenuCategory = "Others"
)

func (menu *MenuCategory) Scan(value interface{}) error {
	*menu = MenuCategory(value.([]byte))
	return nil
}

func (menu MenuCategory) Value() (driver.Value, error) {
	return string(menu), nil
}

type Menu struct {
	ID uint64 					`gorm:"primaryKey"`
	Name string 				`gorm:"column:name" json:"name"`
	Description string 	`gorm:"column:description" json:"description"`
	Type MenuCategory 	`gorm:"type:ENUM('Food', 'Beverage', 'Snack', 'Fruit', 'Grocery', 'Others');column:type" json:"type"`
	COGS uint 					`gorm:"column:cogs" json:"cogs"`
	RetailPrice uint 		`gorm:"column:retail_price;not null" json:"retail_price"`
	WholesalePrice uint `gorm:"column:wholesale_price;not null" json:"wholesale_price"`
	MinOrderQty uint 		`gorm:"column:min_order_qty;default:1;not null" json:"min_order_qty"`
	MaxOrderQty uint 		`gorm:"column:max_order_qty" json:"max_order_qty"`
	PreOrderDays uint 	`gorm:"column:pre_order_days;default:0;not null" json:"pre_order_days"`
	PreOrderHours uint 	`gorm:"column:pre_order_hours;default:0;not null" json:"pre_order_hours"`
	Discount float64 		`gorm:"column:discount;default:0" json:"discount"`
	Image string 				`gorm:"column:image;not null" json:"image"`
	VendorID uint 			`gorm:"column:vendor_id;not null" json:"vendor_id"`
	Vendor Vendor				
	Status string 			`gorm:"column:status;not null" json:"status"`
	CreatedBy string 		`gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}