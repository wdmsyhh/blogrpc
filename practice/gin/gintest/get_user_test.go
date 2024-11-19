package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	reqBody := map[string]interface{}{"id": 11}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "?name=test", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	user := &User{
		DB:          "test",
		RedisClient: NewTestRedis(), // 依赖注入方便写测试代码
	}

	// 业务逻辑
	user.GetUser(ctx)

	// 检查响应
	assert.EqualValues(t, http.StatusOK, w.Code)
	var respBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		panic("Failed to unmarshal response body")
	}
	assert.EqualValues(t, respBody["id"].(float64), 1)

	// 打印响应
	println(w.Body.String())
}
