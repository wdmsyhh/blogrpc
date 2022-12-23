package log

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const (
	HEADER_X_FORWARDED_FOR = "X-Forwarded-For"
)

var Stdout *log.Logger

func init() {
	Stdout = log.New(os.Stdout, "", 0)
}

type AccessLog struct {
	Type              string                 `json:"type"`
	Method            string                 `json:"method"`
	Url               string                 `json:"url"`
	RemoteAddress     string                 `json:"remoteAddr"`
	RemotePort        string                 `json:"remotePort"`
	StatusCode        int                    `json:"responseStatus"`
	ResponseTime      int64                  `json:"responseTime"`
	Referer           string                 `json:"referer"`
	UserAgent         string                 `json:"userAgent"`
	Body              string                 `json:"body"`
	ContentLength     int                    `json:"responseBodySize"`
	RequestId         string                 `json:"reqId"`
	Host              string                 `json:"host,omitempty"`
	HttpXForwardedFor string                 `json:"httpXForwardedFor,omitempty"`
	StartTime         time.Time              `json:"-"`
	EndTime           time.Time              `json:"-"`
	AuthenticatedUser string                 `json:"authenticatedUser,omitempty"`
	Others            map[string]interface{} `json:"others"`
	TenantId          string                 `json:"tenantId"`
}

func NewAccessLog() *AccessLog {
	return &AccessLog{
		Type:      "access",
		StartTime: time.Now(),
		Others:    map[string]interface{}{},
	}
}

func (al *AccessLog) End() time.Duration {
	elapsed := time.Since(al.StartTime)
	al.EndTime = time.Now()
	al.ResponseTime = elapsed.Nanoseconds() / 1e6
	al.Others["startTime"] = al.StartTime
	al.Others["endTime"] = al.EndTime
	return elapsed
}

func (al *AccessLog) Json() []byte {
	return ToJson(al)
}

func ToJson(v interface{}) []byte {
	buf, _ := json.Marshal(v)
	return buf
}
