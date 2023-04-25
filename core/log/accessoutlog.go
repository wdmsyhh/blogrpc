package log

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/cast"
)

const (
	RESPONSE_BODY_RECORD_SWITCH = "ResponseBodyRecordSwitch"
)

type accessOutLog struct {
	Type             string                 `json:"type"`
	RemoteAddr       string                 `json:"remoteAddr"`
	ReqId            string                 `json:"reqId"`
	Host             string                 `json:"host"`
	Scheme           string                 `json:"scheme"`
	Port             int                    `json:"port"`
	Method           string                 `json:"method"`
	Url              string                 `json:"url"`
	Body             string                 `json:"body"`
	ResponseStatus   int                    `json:"responseStatus"`
	ResponseTime     int64                  `json:"responseTime"`
	TimeUnit         string                 `json:"timeUnit"`
	ResponseBodySize int                    `json:"responseBodySize"`
	Others           map[string]interface{} `json:"others"`
	Error            string                 `json:"error"`
	TenantId         string                 `json:"tenantId"`
}

func NewAccessOutLog(reqId, host, scheme, method, path, body, tenantId string) accessOutLog {
	return accessOutLog{
		ReqId:  reqId,
		Host:   host,
		Scheme: scheme,
		Method: method,
		// the Url is actually the path part, like "/v1/members?page=10"
		Url:      path,
		Body:     body,
		TenantId: tenantId,
		Others:   make(map[string]interface{}),
	}
}

// End shall be called right after the request ends, one
// accessOutLog is able to use End() multiple time
func (a *accessOutLog) End(ctx context.Context, responseStatus int, remoteAddr, responseBody string, startedAt time.Time, err error) {
	now := time.Now()
	if !now.After(startedAt) {
		panic("Invalid request start time")
	}
	elapsed := now.Sub(startedAt)
	a.TimeUnit = "ms"
	a.ResponseTime = elapsed.Nanoseconds() / 1000000

	a.ResponseBodySize = len(responseBody)
	a.ResponseStatus = responseStatus

	if IsRecordResponseBody(ctx) {
		a.Others["responseBody"] = responseBody
	}

	// split out the port part of remoteAddr
	addrs := strings.Split(remoteAddr, ":")
	if len(addrs) > 0 {
		a.RemoteAddr = addrs[0]
	}
	if len(addrs) > 1 {
		a.Port = cast.ToInt(addrs[1])
	}

	// avoid user changing the log type
	a.Type = "accessOut"

	if err != nil {
		a.Error = err.Error()
	}

	bytes := ToJson(a)
	Stdout.Printf("%s", bytes)
}

func SwitchOnResponseBodyLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, RESPONSE_BODY_RECORD_SWITCH, true)
}

func IsRecordResponseBody(ctx context.Context) bool {
	record := ctx.Value(RESPONSE_BODY_RECORD_SWITCH)
	return record != nil && record.(bool)
}
