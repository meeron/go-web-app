package products

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"web-app/database"
)

func GetAll(ctx *gin.Context) {
	db, err := database.Connect()
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	defer db.Close()

	products, err := db.Products.Find()
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	result := make([]Product, 0)

	for _, p := range products {
		result = append(result, Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	ctx.JSON(200, result)
}

func Add(ctx *gin.Context) {
	var body struct {
		Name  string
		Price float32
	}

	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	db, err := database.Connect()
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	defer db.Close()

	newEntity, err := db.Products.Add(database.Product{
		Name:  body.Name,
		Price: body.Price,
	})

	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	newProduct := Product{
		Id:    newEntity.ID,
		Name:  newEntity.Name,
		Price: newEntity.Price,
	}

	ctx.JSON(201, newProduct)
}

func Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(404)
		return
	}

	db, err := database.Connect()
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	defer db.Close()

	product, err := db.Products.GetById(id)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	if product == nil {
		ctx.JSON(422, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.JSON(200, Product{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	})
}

func Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(404)
		return
	}

	db, err := database.Connect()
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	defer db.Close()

	exists, err := db.Products.Remove(id)
	if !exists {
		ctx.JSON(422, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.Status(200)
}
