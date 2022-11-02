package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"web-app/web/jwt"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

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

		if !jwt.Validate(token) {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
