package products

type Product struct {
	Id    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type NewProduct struct {
	Name  string  `json:"name" example:"Apple" validate:"required"`
	Price float32 `json:"price" example:"9.99" validate:"gt=0"`
}
