package models

import (
	"time"
	"database/sql/driver"
)

type UserCategory string

const (
	CustomerType UserCategory = "Customer"
	VendorType UserCategory = "Vendor"
	AdminType UserCategory = "Admin"
)

func (user *UserCategory) Scan(value interface{}) error {
	*user = UserCategory(value.([]byte))
	return nil
}

func (user UserCategory) Value() (driver.Value, error) {
	return string(user), nil
}

type User struct {
	ID uint64						`gorm:"primaryKey"`
	Name string 				`gorm:"column:name;not null" json:"name"`
	Email string				`gorm:"column:email;not null" json:"email"`
	Password string	 		`gorm:"column:password;not null" json:"password"`
	Phone string 				`gorm:"column:phone;not null" json:"phone"`
	Type UserCategory		`gorm:"type:ENUM('Customer', 'Vendor', 'Admin');column:type" json:"type"`
	Status string 			`gorm:"column:status;not null" json:"status"`
	CreatedBy string 		`gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time	`gorm:"autoUpdateTime:false" json:"updated_at"`
}