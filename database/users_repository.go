package database

import (
	"errors"
	"gorm.io/gorm"
	"web-app/shared"
)

type IUsersRepository interface {
	GetByEmail(email string) *User
	Add(entity *User)
	Exists(email string) bool
}

type gormUsersRepository struct {
	db *gorm.DB
}

func (g gormUsersRepository) GetByEmail(email string) *User {
	var user User

	err := g.db.First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	shared.PanicOnErr(err)

	return &user
}

func (g gormUsersRepository) Add(entity *User) {
	shared.PanicOnErr(g.db.Create(entity).Error)
}

func (g gormUsersRepository) Exists(email string) bool {
	var count int64
	err := g.db.Table("users").Where("email = ?", email).Count(&count).Error

	shared.PanicOnErr(err)

	return count > 0
}
