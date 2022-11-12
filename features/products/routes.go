package products

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web-app/database"
	"web-app/shared"
	"web-app/web"
)

func getAll(ctx *gin.Context) {
	db := database.DbCtx()

	products := db.Products().Find()

	result := make([]ProductVm, 0)

	for _, p := range products {
		result = append(result, ProductVm{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func add(ctx *gin.Context) {
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

	ctx.JSON(http.StatusCreated, ProductVm{
		Id:    newEntity.ID,
		Name:  newEntity.Name,
		Price: newEntity.Price,
	})
}

func get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	shared.PanicOnErr(err)

	db := database.DbCtx()

	product := db.Products().GetById(id)
	if product == nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.JSON(http.StatusOK, ProductVm{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	})
}

func remove(ctx *gin.Context) {
	id := shared.Unwrap(strconv.Atoi(ctx.Param("id")))

	db := database.DbCtx()

	result := db.Products().Remove(id)
	if !result {
		ctx.JSON(422, gin.H{"errorCode": "NotFound"})
		return
	}

	ctx.Status(http.StatusOK)
}

func ConfigureRoutes(app *gin.Engine) {
	g := app.Group("/products", web.Auth())
	{
		g.GET("", getAll)
		g.POST("", add)
		g.GET(":id", get)
		g.DELETE(":id", remove)
	}
}
