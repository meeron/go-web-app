package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"web-app/shared"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer handlePanic(ctx)

		parts := strings.Split(ctx.GetHeader("Authorization"), " ")

		if len(parts) != 2 {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		schema := parts[0]
		token := parts[1]

		if schema != "Bearer" {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		if !shared.ValidateToken(token) {
			ctx.Status(http.StatusUnauthorized)
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
