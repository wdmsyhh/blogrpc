package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"

	core_util "blogrpc/core/util"
)

type RewriteMiddleware struct{}

func (this *RewriteMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := setRequestIdEnv(c)

		//Set response header X_REQUEST_ID to transport request id
		c.Header(core_util.RequestIDKey, requestId)

		c.Header("Cache-Control", "no-cache")

		// Set request header X_REQUEST_ID, then grpc-gateway can read requestId from header
		c.Request.Header.Set(core_util.RequestIDKey, requestId)
		// TODO: remove
		c.Header("x-req-id", requestId)
		c.Request.Header.Set("x-req-id", requestId)

		c.Next()
	}
}

/**
 * 1. Get request id from request header X-Request-Id,
 * if not exists, generate a new request id using github.com/satori/go.uuid
 * 2. Set request id to rest Request Env
 */
func setRequestIdEnv(c *gin.Context) string {
	// TODO: remove x-req-id
	requestId := c.Request.Header.Get("x-req-id")
	if requestId == "" {
		requestId = c.Request.Header.Get(core_util.RequestIDKey)
	}
	if requestId == "" {
		requestId = uuid.NewV4().String()
	}
	c.Set(core_util.RequestIDKey, requestId)
	c.Set("x-req-id", requestId)
	c.Set(core_util.TracingHeaderKey, c.Request.Header.Get(core_util.TracingHeaderKey))
	return requestId
}
