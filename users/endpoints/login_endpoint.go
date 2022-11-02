package endpoints

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"web-app/database"
	"web-app/shared"
)

func Login(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	shared.PanicOnErr(ctx.BindJSON(&body))

	db := database.DbCtx()

	user := db.Users().GetByEmail(body.Email)
	if user == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	hash := sha256.New()
	_, err := hash.Write([]byte(body.Password))
	shared.PanicOnErr(err)

	hashedPass := fmt.Sprintf("%x", hash.Sum(nil))

	if len(hashedPass) != len(user.Password) {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	for i := 0; i < len(hashedPass); i++ {
		if hashedPass[i] != user.Password[i] {
			ctx.Status(http.StatusUnauthorized)
			return
		}
	}

	accessToken := shared.CreateJwtToken(map[string]string{
		"sub":   fmt.Sprintf("%d", user.ID),
		"email": user.Email,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
