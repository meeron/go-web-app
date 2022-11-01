package database

import (
	"errors"
	"gorm.io/gorm"
)

type IUsersRepository interface {
	GetByEmail(email string) (*User, error)
	Add(entity *User) error
	Exists(email string) (bool, error)
}

type gormUsersRepository struct {
	db *gorm.DB
}

func (g gormUsersRepository) GetByEmail(email string) (*User, error) {
	var user User

	err := g.db.First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (g gormUsersRepository) Add(entity *User) error {
	return g.db.Create(entity).Error
}

func (g gormUsersRepository) Exists(email string) (bool, error) {
	var count int64
	err := g.db.Table("users").Where("email = ?", email).Count(&count).Error

	return count > 0, err
}
