package main

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorize := ctx.GetHeader("Authorization")
		if authorize == "" {
			ctx.Status(401)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
