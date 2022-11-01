package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
)

type Product struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	products := make(map[int]Product, 0)

	app := gin.Default()

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"version": "1.0.0",
		})
	})
	app.GET("/products", func(ctx *gin.Context) {
		result := make([]Product, 0)

		for _, product := range products {
			result = append(result, product)
		}

		ctx.JSON(200, result)
	})
	app.POST("/products", func(ctx *gin.Context) {
		var newProduct Product

		err := ctx.BindJSON(&newProduct)
		if err != nil {
			ctx.JSON(500, err)
			return
		}

		newProduct.Id = rand.Intn(9999)

		products[newProduct.Id] = newProduct

		ctx.JSON(201, newProduct)
	})
	app.GET("/products/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.Status(404)
			return
		}

		product, exists := products[id]
		if !exists {
			ctx.JSON(422, gin.H{"errorCode": "NotFound"})
			return
		}

		ctx.JSON(200, product)
	})
	app.DELETE("/products/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.Status(404)
			return
		}

		_, exists := products[id]
		if !exists {
			ctx.JSON(422, gin.H{"errorCode": "NotFound"})
			return
		}

		delete(products, id)

		ctx.Status(200)
	})

	app.Run()
}
