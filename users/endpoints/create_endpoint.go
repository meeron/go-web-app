package endpoints

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web-app/database"
)

func Create(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
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

	exists, err := db.Users.Exists(body.Email)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	if exists {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errorCode": "Exists"})
		return
	}

	newUser := database.User{
		Email:    body.Email,
		Password: body.Password,
	}

	if err = db.Users.Add(&newUser); err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"Id": newUser.ID})
}
