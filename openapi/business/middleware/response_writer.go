package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"

	"blogrpc/core/log"

	"github.com/gin-gonic/gin"
)

var (
	// 保存数据权限验证未通过的接口，不能直接返回 403 需要处理正常返回空值
	responseMapFor403Routers = map[string]interface{}{}
)

type ResponseWriterMiddleware struct {
	gin.ResponseWriter
	ResponseBody *bytes.Buffer
}

func (rw ResponseWriterMiddleware) Write(data []byte) (int, error) {
	return rw.ResponseBody.Write(data)
}

func (ResponseWriterMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		rw := &ResponseWriterMiddleware{ResponseBody: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = rw
		c.Next()
		statusCode := c.Writer.Status()
		if statusCode >= 500 || log.IsRecordResponseBody(c) {
			c.Keys["responseBody"] = rw.ResponseBody.String()
		} else if statusCode == 403 {
			// 处理数据权限返回 403 状态，但需要正常响应空值的接口
			for r, resp := range responseMapFor403Routers {
				if !isMatched(r, c.Request.URL.Path) {
					continue
				}
				rw.ResponseBody.Reset()
				b, _ := json.Marshal(resp)
				rw.ResponseBody.Write(b)
				c.Status(http.StatusOK)
				break
			}
		}
		rw.ResponseWriter.Write(rw.ResponseBody.Bytes())
	}
}

func isMatched(pattern, s string) bool {
	isMatched, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err)
	}

	return isMatched
}
