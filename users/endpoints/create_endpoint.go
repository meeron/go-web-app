package endpoints

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web-app/database"
	"web-app/shared"
)

func Create(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	shared.PanicOnErr(ctx.BindJSON(&body))

	db := database.Connect()
	defer db.Close()

	if exists := db.Users().Exists(body.Email); exists {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errorCode": "Exists"})
		return
	}

	newUser := database.User{
		Email:    body.Email,
		Password: body.Password,
	}

	db.Users().Add(&newUser)

	ctx.JSON(http.StatusCreated, gin.H{"Id": newUser.ID})
}
