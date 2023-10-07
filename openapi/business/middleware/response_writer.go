package middleware

import (
	"blogrpc/core/util"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"os"
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
		if statusCode == http.StatusOK {
			body := map[string]interface{}{}
			json.Unmarshal(rw.ResponseBody.Bytes(), &body)
			hostname, _ := os.Hostname()
			if v, ok := body["service"]; ok {
				body["service"] = "openapi-business-" + hostname + "-" + util.GetIp() + ";" + cast.ToString(v)
			}
			data, _ := json.Marshal(body)
			rw.ResponseWriter.Write(data)
		} else {
			rw.ResponseWriter.Write(rw.ResponseBody.Bytes())
		}
	}
}
