package models

import (
	"database/sql/driver"
	"time"
)

type AdminStatus string

const (
	Active AdminStatus = "Active"
	Inactive AdminStatus = "Inactive"
)

func (admin *AdminStatus) Scan(value interface{}) error {
	*admin = AdminStatus(value.([]byte))
	return nil
}

func (admin AdminStatus) Value() (driver.Value, error) {
	return string(admin), nil
}

type Admin struct {
	ID uint64						`gorm:"primaryKey" json:"id"`
	UserID uint64				`gorm:"column:user_id;not null"`
	User User						`json:"user"`
	Name string 				`json:"name" gorm:"column:name;not null"`
	Email string				`json:"email" gorm:"column:email;not null"`
	Phone string 				`json:"phone" gorm:"column:phone;not null"`
	Status AdminStatus 	`json:"status" gorm:"type:ENUM('Active', 'Inactive);column:status"`
	CreatedAt time.Time	`gorm:"column:created_at;not null"`
	UpdatedAt time.Time	`gorm:"column:updated_at;autoUpdateTime:false"`
	CreatedBy string		`gorm:"column:created_by;not null"`
}