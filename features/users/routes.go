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

// @Summary Authenticate user
// @Schemes
// @Tags Users
// @Produce json
// @Accept json
// @Param request body users.Login true "Login credentials"
// @Success 200 {object} users.Token
// @Failure 401
// @Router /login [post]
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

	ctx.JSON(http.StatusOK, Token{
		AccessToken: accessToken,
	})
}

// @Summary Create user
// @Schemes
// @Tags Users
// @Produce json
// @Accept json
// @Param request body users.NewUser true "New user data"
// @Success 201 {object} users.User
// @Failure 422 {object} web.Error "Exists"
// @Router /users [post]
func create(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	shared.PanicOnErr(ctx.BindJSON(&body))

	db := database.DbCtx()

	if exists := db.Users().Exists(body.Email); exists {
		ctx.JSON(http.StatusUnprocessableEntity, web.Exists())
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

	ctx.JSON(http.StatusCreated, User{
		Id: int(newUser.ID),
	})
}

func ConfigureRoutes(app *gin.Engine) {
	app.POST("/login", login)

	g := app.Group("/users", web.Auth())
	{
		g.POST("", create)
	}
}
