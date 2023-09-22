package controller

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var testKey = []byte("123")

const EXPIRE_TIME = 60 * 60

type FormGetAccessToken struct {
	AppId     string `json:"appId" form:"appId"`
	AppSecret string `json:"appSecret" form:"appSecret"`
	Aid       string `json:"aid" form:"aid"`
}

func AccessTokenHandler(c *gin.Context) {
	var payload FormGetAccessToken
	err := c.BindQuery(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	appId := payload.AppId
	appSecret := payload.AppSecret
	aid := payload.Aid

	if aid == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "Miss parameter of aid"})
		return
	}
	if appId == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "Miss parameter of app_id"})
		return
	}
	if appSecret == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "Miss parameter of app_secret"})
		return
	}

	token, err := generateAccessToken(appId, appSecret, aid, EXPIRE_TIME)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": token,
		"expires_in":   EXPIRE_TIME,
	})
}

func generateAccessToken(appId, appSecret, aid string, expireTime int64) (string, error) {
	method := jwt.GetSigningMethod("HS256")
	if method == nil {
		return "", errors.New("Failed to get the signing method.")
	}

	token := jwt.New(method)
	iat := time.Now().Unix()
	exp := iat + expireTime
	claims := jwt.MapClaims{
		"iat":       iat,
		"exp":       exp,
		"iss":       "iss",
		"aud":       "aud",
		"sub":       fmt.Sprintf("%s:%s", appId, appSecret),
		"appId":     appId,
		"appSecret": appSecret,
		"aid":       aid,
	}

	header := map[string]interface{}{
		"typ": "JWT",
		"alg": method.Alg(),
	}

	token.Claims = claims
	token.Header = header

	out, err := token.SignedString(testKey)

	if err != nil {
		return "", err
	}

	return out, nil
}
