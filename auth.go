package main

import (
	"github.com/gin-gonic/gin"
)

const fakeUser = "admin"
const fakePass = "admin"

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, pass, hasAuth := ctx.Request.BasicAuth()

		if !hasAuth {
			ctx.Status(401)
			ctx.Abort()
			return
		}

		if user != fakeUser || pass != fakePass {
			ctx.Status(401)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
