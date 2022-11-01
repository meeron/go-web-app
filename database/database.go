package database

type DbContext struct {
	Products IProductsRepository
}

func Connect() (*DbContext, error) {
	return &DbContext{
		Products: &memoryProductsRepository{},
	}, nil
}
