package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func Auth(c *gin.Context) {
	accessToken := c.Request.Header.Get("X-Access-Token")
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "Miss parameter of X-Access-Token"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return testKey, nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message1": err.Error(),
		})
		c.Abort()
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	err = claims.Valid()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message2": err.Error(),
		})
		c.Abort()
		return
	}

	appId := cast.ToString(claims["appId"])
	appSecret := cast.ToString(claims["appSecret"])

	c.Request.Header.Set("App-Id", appId)
	c.Request.Header.Set("App-Secret", appSecret)

	c.Next()
}
