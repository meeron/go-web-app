package database

import (
	"errors"
	"gorm.io/gorm"
	"web-app/shared"
)

type IProductsRepository interface {
	Find() []Product
	Add(new *Product)
	GetById(id int) *Product
	Remove(id int) bool
}

type gormProductsRepository struct {
	db *gorm.DB
}

func (repo gormProductsRepository) Add(new *Product) {
	shared.PanicOnErr(repo.db.Create(new).Error)
}

func (repo gormProductsRepository) GetById(id int) *Product {
	var product Product

	err := repo.db.First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	shared.PanicOnErr(err)

	return &product
}

func (repo gormProductsRepository) Find() []Product {
	products := make([]Product, 0)

	shared.PanicOnErr(repo.db.Find(&products).Error)

	return products
}

func (repo gormProductsRepository) Remove(id int) bool {
	var product Product

	result := repo.db.Delete(&product, id)
	shared.PanicOnErr(result.Error)

	if result.RowsAffected == 0 {
		return false
	}

	return true
}
