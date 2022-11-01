package database

import (
	"errors"
	"gorm.io/gorm"
)

type IProductsRepository interface {
	Find() ([]Product, error)
	Add(new Product) (Product, error)
	GetById(id int) (*Product, error)
	Remove(id int) (bool, error)
}

type gormProductsRepository struct {
	db *gorm.DB
}

func (repo gormProductsRepository) Add(new Product) (Product, error) {
	err := repo.db.Create(&new).Error

	return new, err
}

func (repo gormProductsRepository) GetById(id int) (*Product, error) {
	var product Product

	err := repo.db.First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &product, err
}

func (repo gormProductsRepository) Find() ([]Product, error) {
	products := make([]Product, 0)

	err := repo.db.Find(&products).Error

	return products, err
}

func (repo gormProductsRepository) Remove(id int) (bool, error) {
	var product Product

	result := repo.db.Delete(&product, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, result.Error
}
