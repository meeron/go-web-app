package endpoints

import (
	"crypto/sha256"
	"fmt"
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

	db := database.DbCtx()

	if exists := db.Users().Exists(body.Email); exists {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errorCode": "Exists"})
		return
	}

	hash := sha256.New()
	_, err := hash.Write([]byte(body.Password))
	shared.PanicOnErr(err)

	newUser := database.User{
		Email:    body.Email,
		Password: fmt.Sprintf("%x", hash.Sum(nil)),
	}

	db.Users().Add(&newUser)

	ctx.JSON(http.StatusCreated, gin.H{"Id": newUser.ID})
}
