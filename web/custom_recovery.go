package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CustomRecovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := recover()
		if err == nil {
			ctx.Next()
			return
		}

		// TODO: Log error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": "Generic",
			"message":   err.(error).Error(),
		})
		ctx.Abort()
	}
}
