package users

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"web-app/database"
	"web-app/shared"
	"web-app/web"
	"web-app/web/jwt"
)

func login(ctx *gin.Context) {
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

	accessToken := jwt.Create(map[string]string{
		"sub":   fmt.Sprintf("%d", user.ID),
		"email": user.Email,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func create(ctx *gin.Context) {
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

func ConfigureRoutes(app *gin.Engine) {
	app.POST("/login", login)

	g := app.Group("/users", web.Auth())
	{
		g.POST("", create)
	}
}
