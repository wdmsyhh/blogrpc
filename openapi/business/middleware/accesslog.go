package middleware

import (
	"blogrpc/core/log"
	core_util "blogrpc/core/util"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

type AccessLogJsonMiddleware struct {
}

func (self *AccessLogJsonMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessLog := initAccessLog(c)
		// call the handler
		c.Next()
		path := c.Request.URL.Path
		if !core_util.StrInArray(path, &[]string{"/", "/ping", "/v1/oss/azureBlob"}) {
			recordAccessLog(c, accessLog)
		}
	}
}

func initAccessLog(c *gin.Context) *log.AccessLog {
	body, _ := ioutil.ReadAll(c.Request.Body)
	// Reset c.Request.Body so it can be use again
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	accessLog := log.NewAccessLog()
	accessLog.Body = string(body)
	return accessLog
}

func recordAccessLog(c *gin.Context, accessLog *log.AccessLog) {
	remoteAddr, remotePort, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))

	accessLog.Method = c.Request.Method
	accessLog.Url = c.Request.URL.RequestURI()
	accessLog.RemoteAddress = remoteAddr
	accessLog.RemotePort = remotePort
	accessLog.StatusCode = c.Writer.Status()
	accessLog.Referer = c.Request.Referer()
	accessLog.UserAgent = c.Request.UserAgent()
	accessLog.ContentLength = c.Writer.Size()
	accessLog.RequestId = c.Value(core_util.RequestIDKey).(string)
	accessLog.Others["category"] = "openapi"
	accessLog.Host = c.Request.Host
	accessLog.HttpXForwardedFor = c.Request.Header.Get(log.HEADER_X_FORWARDED_FOR)
	if accessLog.StatusCode == http.StatusUnauthorized {
		// 鉴权失败不输出 body
		accessLog.Body = ""
	}
	if accessLog.StatusCode >= 500 || log.IsRecordResponseBody(c) {
		if responseBody, ok := c.Get("responseBody"); ok {
			accessLog.Others["responseBody"] = responseBody
		}
	}

	accountId, exists := c.Get(core_util.AccountIdKey)
	if exists {
		accessLog.TenantId = cast.ToString(accountId)
	}

	role := c.Request.Header.Get(core_util.AUTHENTICATED_USER_ROLE_IN_HEADER)
	userId := c.Request.Header.Get(core_util.AUTHENTICATED_USER_IN_HEADER)

	if userId != "" {
		if role == "user" {
			role = "portal"
		}
		accessLog.AuthenticatedUser = fmt.Sprintf("business-%s/%s", role, userId)
	}

	elapsed := accessLog.End()

	// record the elapsed time the request processed
	c.Set("ELAPSED_TIME", &elapsed)

	// record accessLog
	log.Stdout.Printf("%s", accessLog.Json())
}
