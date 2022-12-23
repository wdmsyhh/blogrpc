package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	blogrpc_codes "blogrpc/core/codes"
	blogrpc_error "blogrpc/core/errors"
	blogrpc_util "blogrpc/core/util"
	blogrpc_types "blogrpc/proto/common/types"

	"blogrpc/core/extension/bson"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	SECONDS_PER_HOUR = 3600
	SECONDS_PER_DAY  = 86400
	MILLIS_OF_SECOND = 1000

	TIME_REGX = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}([+|-]\\d{2}:\\d{2}|Z)$"

	INVALID            = "invalid"
	MISSING_FIELD      = "missing_field"
	DATE_FORMAT        = "2006-01-02"
	SECRET_HASH_PREFIX = "scrm:jwt:"

	TOKEN_SUB_PREFIX   = "app:"          // the prefix for accessToken's sub field
	TOKEN_EXIST_PREFIX = "access_token:" // used to check if the accessToken is deleted
	TOKEN_MAP_PREFIX   = "api_token_map:"

	// TODO: 合并配置

	TOKEN_AUDIENCE_OPENAPI = "openAPI"
	TOKEN_AUDIENCE_PORTAL  = "portal"
	TOKEN_AUDIENCE_STAFF   = "staff"

	PORTAL_TOKEN_FORMAT = "portal:token:%s"
	STAFF_TOKEN_FORMAT  = "staff:token:%s"
	APP_TOKEN_FORMAT    = "access_token:%s"
)

var (
	rxTimeStr = regexp.MustCompile(TIME_REGX)
)

func Md5(origin string) string {
	h := md5.New()
	io.WriteString(h, origin)
	return hex.EncodeToString(h.Sum(nil))
}

func RandomStr(count int) string {
	if count == 0 {
		return ""
	}
	charList := "0123456789abcdefghigklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	buffer := make([]rune, count)
	for i := 0; i < count; i++ {
		buffer[i] = rune(charList[r.Intn(len(charList))])
	}
	return string(buffer)
}

func RandomNumber(max int, number int) []int {
	buffer := make([]int, number)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < number; i++ {
		buffer[i] = r.Intn(max)
	}
	return buffer
}

func RandomStrFromSpecialStr(count int, charList string) string {
	if count == 0 || charList == "" {
		return ""
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	buffer := make([]rune, count)
	for i := 0; i < count; i++ {
		buffer[i] = rune(charList[r.Intn(len(charList))])
	}
	return string(buffer)
}

func EnsureIntVal(value string, defaultVal int) int {
	retVal := defaultVal
	if value != "" {
		retVal = cast.ToInt(value)
	}
	return retVal
}

func IsEmptyValues(args ...string) bool {
	for _, arg := range args {
		if arg != "" {
			return false
		}
	}
	return true
}

func TransRangeType(left, right string) blogrpc_types.RangeType {
	switch fmt.Sprintf("%s-%s", left, right) {
	case "open-open":
		return blogrpc_types.RangeType_OPEN_OPEN
	case "open-close":
		return blogrpc_types.RangeType_OPEN_CLOSE
	case "open-infinite":
		return blogrpc_types.RangeType_OPEN_INFINITE
	case "close-open":
		return blogrpc_types.RangeType_CLOSE_OPEN
	case "close-close":
		return blogrpc_types.RangeType_CLOSE_CLOSE
	case "close-infinite":
		return blogrpc_types.RangeType_CLOSE_INFINITE
	case "infinite-open":
		return blogrpc_types.RangeType_INFINITE_OPEN
	case "infinite-close":
		return blogrpc_types.RangeType_INFINITE_CLOSE
	}
	return blogrpc_types.RangeType(-1)
}

func TransSecTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp, 0).Format(time.RFC3339)
}

func TransIntTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp/1000, 0).Format(time.RFC3339)
}

func TransStrTimestamp(timeStr string) string {
	t, _ := time.Parse(time.RFC3339Nano, timeStr)
	return TransTime(t)
}

func TransIntTimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}

func TransStrToSecondTimestamp(timeStr string) int64 {
	t, err := time.Parse(time.RFC3339Nano, timeStr)
	if nil != err {
		return 0
	}
	return t.Unix()
}

func TransStrToTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func TransTime(timestamp interface{}) string {
	return timestamp.(time.Time).Format(time.RFC3339)
}

func FirstDayOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

func TodayZeroTime() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func IsTimeStr(str string) bool {
	return rxTimeStr.MatchString(str)
}

func GetStartTimeOfDay(argTime time.Time) time.Time {
	year, month, day := argTime.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// Common HTTP resp
func HttpNoContent(resource interface{}, err error, c *gin.Context) bool {
	if reflect.ValueOf(resource).IsNil() && nil == err {
		c.Abort()
		c.Status(http.StatusNoContent)
		return true
	}

	return false
}

func HttpErrorUnprocessableEntity(c *gin.Context, msg string) {
	HttpError(c, map[string]string{"message": msg}, http.StatusUnprocessableEntity)
}

func HttpError(c *gin.Context, error interface{}, code int) {
	c.JSON(code, error)
}

func Http422(c *gin.Context, err error) bool {
	if nil != err {
		c.Abort()
		HttpErrorUnprocessableEntity(c, err.Error())
		return true
	}

	return false
}

func ParseRPCError(c *gin.Context, err *blogrpc_error.BlogrpcError) bool {
	if err != nil {
		HttpErrorUnprocessableEntity(c, err.Desc)
		return true
	}
	return false
}

func ParseRPCErrorWithErrorMap(c *gin.Context, err *blogrpc_error.BlogrpcError, errorMap map[blogrpc_codes.Code]string) bool {
	if err != nil {
		message := err.Desc
		if msg, ok := errorMap[err.Code]; ok {
			message = msg
		}
		HttpErrorUnprocessableEntity(c, message)
		return true
	}
	return false
}

func StrInArray(search string, items *[]string) bool {
	contains := false
	for _, item := range *items {
		if item == search {
			contains = true
		}
	}
	return contains
}

func IsArray(value interface{}) bool {
	switch value.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}

func IsEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func LowercaseFirst(word string) string {
	length := len(word)
	if length == 0 {
		return ""
	}
	remaining := word[1:]
	first := strings.ToLower(string(word[0]))
	return strings.Join([]string{first, remaining}, "")
}

func EncryptSha256(pData *[]byte, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(*pData)
	bytes := mac.Sum(nil)
	return hex.EncodeToString(bytes)
}

func ToString(data interface{}) string {
	switch data.(type) {
	default:
		return ""
	case bool:
		if data.(bool) {
			return "true"
		} else {
			return "false"
		}
	case int:
		return cast.ToString(data.(int))
	case int32:
		return strconv.FormatInt(int64(data.(int64)), 10)
	case uint32:
		return strconv.FormatInt(int64(data.(uint32)), 10)
	case int64:
		return strconv.FormatInt(int64(data.(int64)), 10)
	case float32:
		return strconv.FormatFloat(float64(data.(float32)), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(data.(float64), 'f', -1, 64)
	case string:
		return data.(string)
	}
}

func ValidateOrigin(origin string) bool {
	result := true
	availableOrigins := []string{"wechat", "weibo", "alipay", "facebook", "portal", "app:android", "app:ios", "app:web", "app:webview", "others"}
	if !StrInArray(origin, &availableOrigins) {
		result = false
	}
	return result
}

func GetVersion() string {
	return "v1"
}

func BuildUniqueErrorMap(message string, resource string, resourceValue string) map[string]interface{} {
	errMap := map[string]string{
		resource: resourceValue,
		"code":   "already_exist",
	}
	return map[string]interface{}{"message": message, "errors": errMap}
}

func BuildRequiredErrorMap(field string, message string, code string) map[string]interface{} {
	errMap := map[string]string{
		field:  message,
		"code": code,
	}
	return map[string]interface{}{"message": message, "errors": errMap}
}

func GetClientIp(request *http.Request) string {
	remoteAddr := request.Header.Get("X-Real-IP")
	if remoteAddr != "" {
		return remoteAddr
	}
	xff := request.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		remoteAddr = strings.TrimSpace(parts[0])
		if remoteAddr != "" {
			return remoteAddr
		}
	}
	remoteAddr = request.RemoteAddr
	if strings.Contains(remoteAddr, ":") {
		ips := strings.Split(remoteAddr, ":") //remote addr is: ip:port format
		remoteAddr = ips[0]                   //only get the ip part, skip the port part
	}
	return remoteAddr
}

func IsPrivateIP(ip string) bool {
	private := false
	IP := net.ParseIP(ip)
	if IP != nil {
		_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
		_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
		_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
		private = private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)
	}
	return private
}

func GetBussinessDomain() string {
	return viper.GetString("business-api-domain")
}

func GetExtensionDomain(serviceName string) string {
	extensionRequest := viper.GetStringMap("extension-request")
	key := fmt.Sprintf("domain-%s", serviceName)
	return cast.ToString(extensionRequest[key])
}

func IsMongoId(checkValue string) bool {
	if !bson.IsObjectIdHex(checkValue) {
		return false
	}
	return true
}

func ExtendMap(dst, src interface{}) {
	dstValue, srcValue := reflect.ValueOf(dst), reflect.ValueOf(src)

	for _, key := range srcValue.MapKeys() {
		dstValue.SetMapIndex(key, srcValue.MapIndex(key))
	}
}

func Round(number float64) uint64 {
	result := math.Ceil(number)
	if (result - number) > 0.50000000001 {
		result -= 1
	}
	return uint64(result)
}

func DecodeMapToStruct(m map[string]interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(config)
	return decoder.Decode(m)
}

func DecodeMapToStructByWeaklyTyped(m map[string]interface{}, output interface{}) error {
	// Be compatible with weakly typed input.
	// Reference: https://pkg.go.dev/github.com/mitchellh/mapstructure#example-Decode-WeaklyTypedInput.
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Metadata:         nil,
		Result:           output,
		TagName:          "json",
	}
	decoder, _ := mapstructure.NewDecoder(config)
	return decoder.Decode(m)
}

func ExceedLimit(limit int, items ...[]string) bool {
	for _, item := range items {
		if len(item) > limit {
			return true
		}
	}

	return false
}

func GetGrpcContext(c *gin.Context) context.Context {
	//Refer: https://godoc.org/google.golang.org/grpc/metadata
	ctx := context.Background()

	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.New(map[string]string{
			blogrpc_util.AccountIdKey:     c.MustGet(blogrpc_util.AccountIdKey).(string),
			blogrpc_util.RequestIDKey:     c.MustGet(blogrpc_util.RequestIDKey).(string),
			blogrpc_util.TracingHeaderKey: c.Request.Header.Get(blogrpc_util.TracingHeaderKey),
		}),
	)

	return ctx
}

func GetGrpcContextWithoutAid(c *gin.Context) context.Context {
	//Refer: https://godoc.org/google.golang.org/grpc/metadata
	ctx := context.Background()

	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.New(map[string]string{
			blogrpc_util.RequestIDKey:     c.MustGet(blogrpc_util.RequestIDKey).(string),
			blogrpc_util.TracingHeaderKey: c.Request.Header.Get(blogrpc_util.TracingHeaderKey),
		}),
	)

	return ctx
}
