package endpoints

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(ctx *gin.Context) {
	ctx.Status(http.StatusUnauthorized)
}
