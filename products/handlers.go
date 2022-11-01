package products

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
)

func GetAll(ctx *gin.Context) {
	result := make([]Product, 0)

	//for _, product := range products {
	//	result = append(result, product)
	//}

	ctx.JSON(200, result)
}

func Add(ctx *gin.Context) {
	var newProduct Product

	err := ctx.BindJSON(&newProduct)
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	newProduct.Id = rand.Intn(9999)

	//products[newProduct.Id] = newProduct

	ctx.JSON(201, newProduct)
}

func Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(404)
		return
	}

	//product, exists := products[id]
	//if !exists {
	//	ctx.JSON(422, gin.H{"errorCode": "NotFound"})
	//	return
	//}

	//ctx.JSON(200, product)
	ctx.JSON(200, id)
}

func Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(404)
		return
	}

	/*
		_, exists := products[id]
		if !exists {
			ctx.JSON(422, gin.H{"errorCode": "NotFound"})
			return
		}

		delete(products, id)
	*/

	ctx.JSON(200, id)
}
