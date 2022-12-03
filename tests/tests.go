package tests

import (
	"database/sql/driver"
	"time"

	"github.com/gin-gonic/gin"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	return router
}
