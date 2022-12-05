package users

import (
	"crypto/sha256"
	"fmt"
	"web-app/database"
	"web-app/shared"
	"web-app/web/jwt"

	"github.com/gofiber/fiber/v2"
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
func login(ctx *fiber.Ctx) error {

	var body Login

	parserErr := ctx.BodyParser(&body)
	if parserErr != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DbCtx()

	user := db.Users().GetByEmail(body.Email)
	if user == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	hash := sha256.New()
	_, err := hash.Write([]byte(body.Password))
	shared.PanicOnErr(err)

	hashedPass := fmt.Sprintf("%x", hash.Sum(nil))

	if len(hashedPass) != len(user.Password) {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	for i := 0; i < len(hashedPass); i++ {
		if hashedPass[i] != user.Password[i] {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
	}

	accessToken := jwt.Create(map[string]string{
		"sub":   fmt.Sprintf("%d", user.ID),
		"email": user.Email,
	})

	return ctx.JSON(Token{
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
func create() {
	/*
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
	*/
}

func ConfigureRoutes(app *fiber.App) {
	app.Post("/login", login)

	/*
		g := app.Group("/users", web.Auth())
		{
			g.POST("", create)
		}
	*/
}
