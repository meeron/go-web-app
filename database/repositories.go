package database

import (
	"math/rand"
)

var products = make(map[int]Product, 0)

type IProductsRepository interface {
	Find() ([]Product, error)
	Add(new Product) (Product, error)
	GetById(id int) (*Product, error)
	Remove(id int) (bool, error)
}

type memoryProductsRepository struct {
}

func (repo memoryProductsRepository) Add(new Product) (Product, error) {
	new.Id = rand.Intn(999999)

	products[new.Id] = new

	return new, nil
}

func (repo memoryProductsRepository) GetById(id int) (*Product, error) {
	product, exists := products[id]
	if !exists {
		return nil, nil
	}

	return &product, nil
}

func (repo memoryProductsRepository) Find() ([]Product, error) {
	result := make([]Product, 0)

	for _, product := range products {
		result = append(result, product)
	}

	return result, nil
}

func (repo memoryProductsRepository) Remove(id int) (bool, error) {
	_, exists := products[id]
	if !exists {
		return false, nil
	}

	delete(products, id)

	return true, nil
}
