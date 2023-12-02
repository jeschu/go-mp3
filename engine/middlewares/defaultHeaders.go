package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"runtime"
)

const XRequestId = "X-Request-Id"

func DefaultHeaders(server string, requestId bool) gin.HandlerFunc {
	if server == "" {
		server = "go-mp3 " + runtime.Version()
	}
	return func(c *gin.Context) {
		c.Header("Server", server)
		if requestId {
			rqId := c.GetHeader(XRequestId)
			if rqId == "" {
				rqId = uuid.New().String()
			}
			c.Header(XRequestId, rqId)
			c.Set(XRequestId, rqId)
		}
	}
}
