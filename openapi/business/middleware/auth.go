package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"blogrpc/core/extension"
	"blogrpc/core/log"
	core_util "blogrpc/core/util"
	rpc_service "blogrpc/openapi/business/rpc"
	"blogrpc/openapi/business/util"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/cast"
	conf "github.com/spf13/viper"
)

const (
	AUTH_ORDINARY_PATTERN = "^(\\/%s|\\/%s)|modules|apps\\/"
)

var (
	TEST_ENV_ACTIONS = []string{"^DELETE/v2/members/[^\\s]"}

	NO_AUTH_ACTIONS = []string{
		"/v2/user/sendResetPasswordEmail",
		"/v2/user/resetPassword",
		"/v2/validateCode",
		"/v2/user/activate",
		"/v2/user/accounts",
		"/v2/login",
		"/v1/jingdong/auth",
	}

	SessionValidPeriodMap = sync.Map{}
)

type AuthMiddleware struct {
	SigningAlgorithm      string
	Version               string
	Env                   string
	LookupFunc            jwt.Keyfunc
	SkipReferrerValidator bool
}

type LoginUserInfo struct {
	UserId   string `json:"userId"`
	TenantId int64  `json:"tenantId"`
	UserName string `json:"userName"`
	Mobile   string `json:"mobile"`
}

func (am *AuthMiddleware) Auth() gin.HandlerFunc {
	if am.SigningAlgorithm == "" {
		am.SigningAlgorithm = "HS256"
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if util.StrInArray(path, &NO_AUTH_ACTIONS) {
			return
		}

		pattern := fmt.Sprintf(AUTH_ORDINARY_PATTERN, am.Version, "v2")
		needAuth, _ := regexp.MatchString(pattern, path)
		if needAuth {
			accountId, token, err := am.validateOauth(c)
			if err != nil {
				if token == nil {
					unauthorized(c)
					return
				} else {
					if err == redis.ErrNil || strings.Contains(err.Error(), "not exists") {
						unauthorized(c)
						return
					}
					forbidden(c, map[string]string{
						"message": err.Error(),
					})
					return
				}
			}

			c.Request.Header.Set(core_util.ACCOUNT_ID_IN_HEADER, accountId)
			c.Request.Header.Set(core_util.ACCESS_TOKEN_IN_HEADER, token.Raw)

			role, id := parseSub(token)
			c.Request.Header.Set(core_util.ROLE_IN_HEADER, role)
			setAuthenticatedKeyInHeader(c, role, id)
			if role == "staff" {
				formatStaffRequest(c, id)
			}
			if role == "user" {
			}
			if !am.isTestEnv() && am.isTestEnvAction(c) {
				forbidden(c, map[string]string{
					"message": "Access to not allowed action",
				})
				return
			}
		}

		c.Next()
	}
}

func setAuthenticatedKeyInHeader(c *gin.Context, role, userId string) {
	c.Request.Header.Set(core_util.AUTHENTICATED_USER_ROLE_IN_HEADER, role)
	c.Request.Header.Set(core_util.AUTHENTICATED_USER_IN_HEADER, userId)
}

func (am *AuthMiddleware) isTestEnv() bool {
	return !strings.Contains(am.Env, "production")
}

func (am *AuthMiddleware) isTestEnvAction(c *gin.Context) bool {
	httpMethod := strings.ToUpper(c.Request.Method)
	path := c.Request.URL.Path

	for _, action := range TEST_ENV_ACTIONS {
		if match, _ := regexp.MatchString(action, fmt.Sprintf("%s%s", httpMethod, path)); match {
			return true
		}
	}

	return false
}

func getAccessToken(req *http.Request) string {
	if accessToken := req.URL.Query().Get("access_token"); accessToken != "" {
		return accessToken
	}

	if authHeader := req.Header.Get("Authorization"); authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if parts[0] == "token" {
			return parts[1]
		}
	}

	// in order to keep consistence of consumerapi, we'll try to get token
	// from "X-Access-Token"
	if authHeader := req.Header.Get(core_util.ACCESS_TOKEN_IN_HEADER); authHeader != "" {
		return authHeader
	}

	env := conf.GetString("env")
	cookieName := "accesstoken"
	if env != "production" {
		cookieName = fmt.Sprintf("%s_%s", strings.Replace(env, "-", "", -1), cookieName)
	}
	if cookie, err := req.Cookie(cookieName); err == nil {
		return core_util.DecodeAccessToken(cookie.Value)
	}

	return ""
}

func (am *AuthMiddleware) validateOauth(c *gin.Context) (string, *jwt.Token, error) {
	accessToken := getAccessToken(c.Request)
	if accessToken == "" {
		return "", nil, errors.New("Get accessToken failed")
	}

	token, err := parseAccessToken(c, accessToken, am.LookupFunc)
	if err != nil {
		log.Warn(c, "Parse accessToken failed", log.Fields{
			"accessToken": accessToken,
		})
		return "", nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if _, isOldToken := claims["uid"]; isOldToken {
		token, err = getNewToken(c, token, am.LookupFunc)
		if err != nil {
			return "", nil, err
		}

		claims = token.Claims.(jwt.MapClaims)
	}

	sub := cast.ToString(claims["sub"])
	accountId := cast.ToString(claims["aid"])
	if accountId == "" {
		return "", token, errors.New("AccountId does not exist")
	}
	if isAccountDisabled(accountId) {
		return "", token, errors.New("Account disabled")
	}
	if claims["aud"] == util.TOKEN_AUDIENCE_PORTAL {
		if xAccountId := getXAccountId(c); xAccountId != "" && xAccountId != accountId {
			return "", nil, errors.New("Invalid accountId")
		}
	}
	c.Set(core_util.AccountIdKey, accountId)
	err = checkToken(c, claims, accessToken)
	if err != nil {
		log.Warn(c, err.Error(), log.Fields{
			"sub": sub,
		})
		return "", token, err
	}

	if !isForOpenAPI(cast.ToString(claims["aud"])) {
		err := errors.New("The token is not for openAPI")
		log.Warn(c, err.Error(), log.Fields{
			"token": token.Raw,
		})
		return "", nil, err
	}

	return cast.ToString(claims["aid"]), token, nil
}

// 根据 accessToken 的类型对其进行验证，目前只有 app 类型的 accessToken 会
// 验证 ip 白名单，其它类型都只验证是否存在
func checkToken(c *gin.Context, claims jwt.MapClaims, accessToken string) error {
	sub := cast.ToString(claims["sub"])
	splits := strings.Split(sub, ":")
	if len(splits) != 2 {
		return errors.New("Invalid accessToken")
	}

	key := ""
	switch splits[0] {
	case "app":
		key = fmt.Sprintf(util.APP_TOKEN_FORMAT, splits[1])
		value, err := extension.RedisClient.Get(key)
		if err != nil {
			return err
		}

		err = validateClientIp(c, value)
		if value != "" && value != "1" && err == nil {
			return nil
		}

		return err
	case "user":
		key = fmt.Sprintf(util.PORTAL_TOKEN_FORMAT, splits[1])
	case "staff":
		key = fmt.Sprintf(util.STAFF_TOKEN_FORMAT, accessToken)
	default:
		return errors.New("No such token type")
	}

	exists, err := extension.RedisClient.Exists(key)
	if err != nil || !exists {
		return errors.New("AccessToken not exists")
	}
	if splits[0] == "user" {
		accountId := cast.ToString(claims["aid"])
		refreshTTL(c, accountId, key)
	}

	return nil
}

func refreshTTL(c *gin.Context, accountId, key string) {
	value, err := extension.RedisClient.Get(key)
	if err != nil {
		return
	}
	cacheSessionValidPeriodTime, ok := SessionValidPeriodMap.Load(accountId)
	if ok {
		extension.RedisClient.SetEx(key, cacheSessionValidPeriodTime.(int64), value)
		return
	}
	resp, err := rpc_service.CAccount.GetSessionValidPeriod(c)
	if err != nil {
		return
	}
	sessionValidPeriodTime := cast.ToInt64(resp.SessionValidPeriod * 60 * 60) // 小时转换成秒
	if sessionValidPeriodTime != 0 {
		extension.RedisClient.SetEx(key, sessionValidPeriodTime, value)
		SessionValidPeriodMap.Store(accountId, sessionValidPeriodTime)
		time.AfterFunc(10*time.Minute, func() {
			SessionValidPeriodMap.Delete(accountId)
		})
	} else {
		extension.RedisClient.SetEx(key, core_util.AccessTokenPortalExpireTime, value)
	}
}

func validateClientIp(c *gin.Context, ips string) error {
	if ips == "" || ips == "1" {
		return nil
	}

	clientIP := util.GetClientIp(c.Request)
	ipWhiteList := strings.Split(ips, ",")
	for _, ip := range ipWhiteList {
		if strings.Contains(ip, "/") {
			_, ipNet, err := net.ParseCIDR(ip)
			if err != nil {
				continue
			}

			ipSplit := strings.Split(ip, "/")
			ipSuffix := ipSplit[1]
			netClientIP, _, err := net.ParseCIDR(clientIP + "/" + ipSuffix)
			if err != nil {
				continue
			}

			if ipNet.Contains(netClientIP) {
				return nil
			}
		} else {
			if clientIP == ip {
				return nil
			}
		}
	}

	return errors.New(fmt.Sprintf("IP %s isn't in allowlist", clientIP))
}

func parseSub(token *jwt.Token) (role, userId string) {
	claims := token.Claims.(jwt.MapClaims)
	sub := strings.Split(cast.ToString(claims["sub"]), ":")
	if len(sub) != 2 {
		return
	}

	role = sub[0]
	userId = sub[1]
	return
}

func parseAccessToken(c *gin.Context, accessToken string, lookup jwt.Keyfunc) (*jwt.Token, error) {
	token, err := jwt.Parse(accessToken, lookup)
	if err != nil {
		log.Warn(c, "Fail to parse oauth token", log.Fields{
			"token": accessToken,
			"error": err,
		})

		return nil, errors.New("Auth header invalid")
	}

	return token, nil
}

func getNewToken(c *gin.Context, token *jwt.Token, lookup jwt.Keyfunc) (*jwt.Token, error) {
	newToken, err := extension.RedisClient.Get(util.TOKEN_MAP_PREFIX + token.Signature)
	if err != nil {
		log.Warn(c, "Fail to get new token through old token", log.Fields{
			"token": token.Raw,
			"error": err,
		})

		return nil, errors.New("Auth header invalid")
	}

	token, err = parseAccessToken(c, newToken, lookup)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func isForOpenAPI(audience string) bool {
	return audience == util.TOKEN_AUDIENCE_OPENAPI || audience == util.TOKEN_AUDIENCE_PORTAL || audience == util.TOKEN_AUDIENCE_STAFF
}

func unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, map[string]string{"message": "Bad credentials"})
	c.Abort()
}

func forbidden(c *gin.Context, extraFields map[string]string) {
	c.JSON(http.StatusForbidden, extraFields)
	c.Abort()
}

func formatStaffRequest(c *gin.Context, staffId string) {
	storeId, err := extension.RedisClient.Hget(fmt.Sprintf(util.STAFF_TOKEN_FORMAT, getAccessToken(c.Request)), "storeId")
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		c.Abort()
		return
	}
	r := c.Request
	q := r.URL.Query()
	bodyMap := map[string]interface{}{}
	bodyContent, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(bodyContent, &bodyMap)

	isUspEnv := strings.Contains(conf.GetString("env"), "usp")
	if r.URL.Query().Get("_ignoreStore") == "" || isUspEnv {
		// write into query
		q.Set("storeId", storeId)
		// write into body
		bodyMap["storeId"] = storeId
	}

	if r.URL.Query().Get("_ignoreStaff") == "" || isUspEnv {
		// write into query
		q.Set("staffId", staffId)
		// write into body
		bodyMap["staffId"] = staffId
	}

	r.URL.RawQuery = q.Encode()
	newBodyContent, _ := json.Marshal(bodyMap)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(newBodyContent))
	r.ContentLength = int64(len(newBodyContent))
}

func getRealIp(c *gin.Context) string {
	if ips := c.Request.Header["X-Real-Ip"]; len(ips) > 0 {
		return strings.TrimSpace(ips[0])
	}
	return ""
}

// 需要确认 *gin.Context 中已经写入 accountId 之后才可使用
func getAccountId(c *gin.Context) string {
	return c.MustGet(core_util.AccountIdKey).(string)
}

func getXAccountId(c *gin.Context) string {
	accountIds, ok := c.Request.Header["X-Account-Id"]
	if !ok || len(accountIds) == 0 {
		return ""
	}
	return accountIds[0]
}

func isAccountDisabled(accountId string) bool {
	status, _ := extension.RedisClient.Hget(fmt.Sprintf("%s:info", accountId), "status")
	return status == "disabled"
}
