package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const fakeUser = "admin"
const fakePass = "admin"

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer handlePanic(ctx)

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

func handlePanic(ctx *gin.Context) {
	if err := recover(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": "Generic",
			"message":   err.(error).Error(),
		})
		ctx.Abort()
	}
}
