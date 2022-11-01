package endpoints

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Create(ctx *gin.Context) {
	ctx.Status(http.StatusMethodNotAllowed)
}
