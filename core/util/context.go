package util

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	RequestIDKey      = "x-request-id"
	AccountIdKey      = "aid"
	TracingHeaderKey  = "sw8"
	DefaultRequestID  = "00000000-0000-0000-0000-000000000000"
	AccountIdRequired = "AccountId Required"

	ACCESS_LOG_EXTRA_KEY = "access-log-extra"

	ACCOUNT_ID_IN_HEADER               = "X-Account-Id"
	ACCESS_TOKEN_IN_HEADER             = "X-Access-Token"
	MEMBER_ID_IN_HEADER                = "X-Member-Id"
	OPEN_ID_IN_HEADER                  = "X-Open-Id"
	UNION_ID_IN_HEADER                 = "X-Union-Id"
	CHANNEL_ID_IN_HEADER               = "X-Channel-Id"
	CHANNEL_ENCRYPTION_KEY_IN_HEADER   = "X-Channel-Encryption-Key"
	CHANNEL_NAME_IN_HEADER             = "X-Channel-NAME"
	CHANNEL_INTEGRATION_MODE_IN_HEADER = "X-Channel-Integration-Mode"
	AUTHORIZE_IN_HEADER                = "X-Authorize"
	ROLE_IN_HEADER                     = "X-Role"
	REQUEST_ID_IN_HEADER               = "X-Request-Id"
	AUTHENTICATED_USER_IN_HEADER       = "X-Authenticated-User"
	AUTHENTICATED_USER_ROLE_IN_HEADER  = "X-Authenticated-User-Role"
	DISTRIBUTOR_ID_IN_HEADER           = "X-Authenticated-Distributor-Id"
	DISTRIBUTOR_IDS_IN_HEADER          = "X-Authenticated-Distributor-Ids"
	DISTRIBUTOR_SELECT_FIELD_IN_HEADER = "X-Authenticated-Distributor-Select-Field"
	X_REQUESTED_WITH_IN_HEADER         = "X-Requested-With"
	IS_NEW_CHANNEL_IN_HEADER           = "X-Is-New-Channel"
	WEBHOOK_SIGNATURE_IN_HEADER        = "X-Webhook-Signature"
)

var (
	TransactionSessionContextType = reflect.TypeOf(new(mongo.SessionContext)).Elem()
)

func CtxWithRequestID(ctx context.Context, rid string) context.Context {
	return setCtxData(ctx, RequestIDKey, rid)
}

func setCtxData(ctx context.Context, key, value string) context.Context {
	ctx = context.WithValue(ctx, key, value)
	md, _ := metadata.FromIncomingContext(ctx)
	if md.Len() != 0 {
		md.Set(key, value)
	} else {
		md = metadata.New(map[string]string{key: value})
	}

	return metadata.NewIncomingContext(ctx, md)
}

func CtxWithAccountID(ctx context.Context, aid string) context.Context {
	return setCtxData(ctx, AccountIdKey, aid)
}

func CtxWithUserID(ctx context.Context, userId string) context.Context {
	return setCtxData(ctx, strings.ToLower(AUTHENTICATED_USER_IN_HEADER), userId)
}

func ExtractRequestIDFromCtx(ctx context.Context) string {
	// TODO: remove x-req-id
	rid := ExtractValueFromCtx(ctx, "x-req-id")
	if rid == "" {
		rid = ExtractValueFromCtx(ctx, RequestIDKey)
	}
	if rid == "" {
		rid = DefaultRequestID
	}
	return rid
}

func ExtractTracingHeaderFromCtx(ctx context.Context) string {
	return ExtractValueFromCtx(ctx, TracingHeaderKey)
}

func ExtractValueFromCtx(ctx context.Context, key string) string {
	var value string
	if ctx == nil || reflect.ValueOf(ctx).IsNil() {
		return value
	}

	// value from grpc
	if value = GetValueFromContext(ctx, key); value != "" {
		return value
	}

	// value from openapi(*gin.Context)
	if value, ok := ctx.Value(key).(string); ok {
		return value
	}

	return value
}

func GetValueFromContext(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md, ok = metadata.FromOutgoingContext(ctx)
	}

	if !ok { // neither can we get md from incoming nor from outgoing
		return ""
	}

	if md[key] == nil || len(md[key]) == 0 {
		return ""
	}

	return md[key][0]
}

// Must get accountId, if accountId is empty, panic
func MustGetAccountId(ctx context.Context) string {
	accountId := ExtractValueFromCtx(ctx, AccountIdKey)
	if accountId != "" {
		return accountId
	}
	panic(fmt.Sprintf("%v,%#v", AccountIdRequired, GetAccountId(ctx)))
}

func GetAccountId(ctx context.Context) string {
	return ExtractValueFromCtx(ctx, AccountIdKey)
}

func GetUserRole(ctx context.Context) string {
	return ExtractValueFromCtx(ctx, strings.ToLower(AUTHENTICATED_USER_ROLE_IN_HEADER))
}

func GetUserId(ctx context.Context) string {
	return ExtractValueFromCtx(ctx, strings.ToLower(AUTHENTICATED_USER_IN_HEADER))
}

func ClientIP(ctx context.Context) (ip string) {
	p, ok := peer.FromContext(ctx)
	if ok {
		remoteAddr := p.Addr.String()
		if strings.Contains(remoteAddr, ":") {
			ips := strings.Split(remoteAddr, ":") //remote addr is: ip:port format
			ip = ips[0]                           //only get the ip part, skip the port part
		}
	}

	return
}

func ClientPort(ctx context.Context) (port string) {
	p, ok := peer.FromContext(ctx)
	if ok {
		remoteAddr := p.Addr.String()
		if strings.Contains(remoteAddr, ":") {
			ips := strings.Split(remoteAddr, ":") //remote addr is: ip:port format
			port = ips[1]
		}
	}

	return
}

func GetGinContextFromContext(ctx context.Context) *gin.Context {
	c := gin.Context{}
	c.Set(AccountIdKey, ExtractValueFromCtx(ctx, AccountIdKey))
	c.Set(RequestIDKey, ExtractValueFromCtx(ctx, RequestIDKey))

	return &c
}

func DuplicateContext(ctx context.Context) context.Context {
	if IsTransactionSessionContext(ctx) {
		return ctx
	}

	duplicateCtx := context.Background()
	if reflect.TypeOf(ctx) == reflect.TypeOf(&gin.Context{}) {
		if ctx.Value(AccountIdKey) != nil {
			duplicateCtx = setCtxData(duplicateCtx, AccountIdKey, ctx.Value(AccountIdKey).(string))
		}
		if ctx.Value(RequestIDKey) != nil {
			duplicateCtx = setCtxData(duplicateCtx, RequestIDKey, ctx.Value(RequestIDKey).(string))
		}
		if ctx.Value(TracingHeaderKey) != nil {
			duplicateCtx = setCtxData(duplicateCtx, TracingHeaderKey, ctx.Value(TracingHeaderKey).(string))
		}
		if ctx.Value(ACCESS_LOG_EXTRA_KEY) != nil {
			duplicateCtx = context.WithValue(duplicateCtx, ACCESS_LOG_EXTRA_KEY, ctx.Value(ACCESS_LOG_EXTRA_KEY))
		}
		return duplicateCtx
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		duplicateCtx = metadata.NewIncomingContext(duplicateCtx, md)
	}
	md, ok = metadata.FromOutgoingContext(ctx)
	if ok {
		duplicateCtx = metadata.NewOutgoingContext(duplicateCtx, md)
	}
	return duplicateCtx
}

// 复制 context 并设置 accountId
func DuplicateContextWithAid(ctx context.Context, accountId string) context.Context {
	if IsTransactionSessionContext(ctx) {
		return ctx
	}

	duplicateCtx := context.Background()
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md["aid"] = []string{accountId}
		duplicateCtx = metadata.NewIncomingContext(duplicateCtx, md)
	}
	md, ok = metadata.FromOutgoingContext(ctx)
	if ok {
		duplicateCtx = metadata.NewOutgoingContext(duplicateCtx, md)
	}
	return duplicateCtx
}

// 复制 context 并移除 accountId
func DuplicateContextWithoutAid(ctx context.Context) context.Context {
	if IsTransactionSessionContext(ctx) {
		return ctx
	}

	duplicateCtx := context.Background()
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if md["aid"] != nil && len(md["aid"]) > 0 {
			md["aid"] = []string{}
		}
		duplicateCtx = metadata.NewIncomingContext(duplicateCtx, md)
	}
	md, ok = metadata.FromOutgoingContext(ctx)
	if ok {
		duplicateCtx = metadata.NewOutgoingContext(duplicateCtx, md)
	}
	return duplicateCtx
}

func IsTransactionSessionContext(ctx context.Context) bool {
	return reflect.TypeOf(ctx).Implements(TransactionSessionContextType)
}

func WriteAccessLogExtra(ctx context.Context, key string, value interface{}) bool {
	v := ctx.Value(ACCESS_LOG_EXTRA_KEY)
	if v == nil {
		return false
	}

	extra, ok := v.(*sync.Map)
	if !ok || extra == nil {
		return false
	}

	extra.Store(key, value)
	return true
}
