package products

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"web-app/database"
)

func GetAll(ctx *gin.Context) {
	db, err := database.Connect()
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	products, err := db.Products.Find()
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, products)
}

func Add(ctx *gin.Context) {
	var body struct {
		Name string
	}

	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	newEntity, err := db.Products.Add(database.Product{
		Name: body.Name,
	})

	if err != nil {
		ctx.JSON(500, err)
		return
	}

	newProduct := Product{
		Id:   newEntity.Id,
		Name: newEntity.Name,
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
		ctx.JSON(500, err)
		return
	}

	product, err := db.Products.GetById(id)
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	if product == nil {
		ctx.JSON(422, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.JSON(200, Product{
		Id:   product.Id,
		Name: product.Name,
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
		ctx.JSON(500, err)
		return
	}

	exists, err := db.Products.Remove(id)
	if !exists {
		ctx.JSON(422, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.Status(200)
}
