package util

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	AccessTokenDefaultAlg = "HS256"

	AccessTokenServiceMobile    = "mobile"
	AccessTokenServicePortal    = "portal"
	AccessTokenMobileExpireTime = 6 * 60 * 60  // 6 hours
	AccessTokenPortalExpireTime = 12 * 60 * 60 // 12 hours

	TokenUserTypeUser     = "user"
	TokenUserTypeMember   = "member"
	TokenUserTypeFollower = "follower"
)

var (
	supportedUserTypes = []string{TokenUserTypeUser, TokenUserTypeMember, TokenUserTypeFollower}
)

func GeneratePortalAccessToken(accountId, secretKey, userId string) (string, error) {
	return generateAccessToken(accountId, secretKey, TokenUserTypeUser, userId, AccessTokenServicePortal, AccessTokenPortalExpireTime)
}

// GenerateAccessToken generate jwt token
func generateAccessToken(accountId, secretKey, userType, userId, service string, expireTime int64) (string, error) {
	if !StrInArray(userType, &supportedUserTypes) {
		return "", errors.New("unsupported user types.")
	}
	// get the signing method
	method := jwt.GetSigningMethod(AccessTokenDefaultAlg)
	if method == nil {
		return "", errors.New("Failed to get the signing method.")
	}

	// create a new token
	token := jwt.New(method)
	iat := time.Now().Unix()
	exp := iat + expireTime
	claims := jwt.MapClaims{
		"iat": iat,
		"exp": exp,
		"iss": service,
		"aud": service,
		"sub": fmt.Sprintf("%s:%s", userType, userId),
		"aid": accountId,
	}

	header := map[string]interface{}{
		"typ": "JWT",
		"alg": method.Alg(),
	}

	token.Claims = claims
	token.Header = header
	out, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return out, nil

}
