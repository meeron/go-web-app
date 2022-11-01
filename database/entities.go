package database

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name  string
	Price float32
}

type User struct {
	gorm.Model
	Email    string
	Password string
}
