package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const VERSION = "0.0.1"

type VersionController struct{}

func (controller VersionController) RegisterRoutes(engine *gin.Engine) {
	engine.GET(
		"/version",
		func(c *gin.Context) {
			c.PureJSON(http.StatusOK, gin.H{"version": VERSION})
		},
	)
}
