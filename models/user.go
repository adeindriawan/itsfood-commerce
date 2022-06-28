package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID 				uint64	`json:"id"`
	Username 	string 	`json:"username"`
	Password	string 	`json:"password"`
}