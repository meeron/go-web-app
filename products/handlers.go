package products

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web-app/database"
	"web-app/shared"
)

func GetAll(ctx *gin.Context) {
	db := database.DbCtx()

	products := db.Products().Find()

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

	db := database.DbCtx()

	newEntity := database.Product{
		Name:  body.Name,
		Price: body.Price,
	}

	db.Products().Add(&newEntity)

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

	db := database.DbCtx()

	product := db.Products().GetById(id)
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
	id := shared.Unwrap(strconv.Atoi(ctx.Param("id")))

	db := database.DbCtx()

	result := db.Products().Remove(id)
	if !result {
		ctx.JSON(422, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.Status(http.StatusOK)
}
