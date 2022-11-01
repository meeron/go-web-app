package products

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web-app/database"
	"web-app/shared"
)

func GetAll(ctx *gin.Context) {
	db, err := database.Connect()
	shared.PanicOnErr(err)

	defer db.Close()

	products, err := db.Products.Find()
	shared.PanicOnErr(err)

	result := make([]Product, 0)

	for _, p := range products {
		result = append(result, Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func Add(ctx *gin.Context) {
	var body struct {
		Name  string
		Price float32
	}

	shared.PanicOnErr(ctx.BindJSON(&body))

	db, err := database.Connect()
	shared.PanicOnErr(err)

	defer db.Close()

	newEntity, err := db.Products.Add(database.Product{
		Name:  body.Name,
		Price: body.Price,
	})
	shared.PanicOnErr(err)

	newProduct := Product{
		Id:    newEntity.ID,
		Name:  newEntity.Name,
		Price: newEntity.Price,
	}

	ctx.JSON(http.StatusCreated, newProduct)
}

func Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	shared.PanicOnErr(err)

	db, err := database.Connect()
	shared.PanicOnErr(err)

	defer db.Close()

	product, err := db.Products.GetById(id)
	shared.PanicOnErr(err)

	if product == nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.JSON(http.StatusOK, Product{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	})
}

func Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	shared.PanicOnErr(err)

	db, err := database.Connect()
	shared.PanicOnErr(err)

	defer db.Close()

	exists, err := db.Products.Remove(id)
	shared.PanicOnErr(err)

	if !exists {
		ctx.JSON(422, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.Status(http.StatusOK)
}
