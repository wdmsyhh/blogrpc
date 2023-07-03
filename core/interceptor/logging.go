package interceptor

import (
	"blogrpc/core/log"
	"blogrpc/core/util"
	"encoding/json"
	"runtime"
	"sync"
	"time"

	"github.com/spf13/cast"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func init() {
	AddInterceptor(&LoggingInterceptor{})
}

type LoggingInterceptor struct {
	env     string
	service string
	address string
}

func (this *LoggingInterceptor) Name() string { return "Logging" }
func (this *LoggingInterceptor) InitWithConf(conf map[string]interface{}, debug bool) error {
	this.service = cast.ToString(conf["service"])
	this.address = cast.ToString(conf["addr"])
	this.env = cast.ToString(conf["env"])
	return nil
}

// Logging interceptor for grpc
func (this *LoggingInterceptor) Handle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	accessLog := initAccessLog(req, info)

	sm := &sync.Map{}
	ctx = context.WithValue(ctx, util.ACCESS_LOG_EXTRA_KEY, sm)
	resp, err = handler(ctx, req)
	// Don't print any log for health checking
	if info.FullMethod == "/Health/Check" {
		return resp, err
	}

	if this.env == "local" {
		accessLog.Others["response"] = resp
	}

	recordAccessLog(ctx, accessLog, err, map[string]string{
		"env":     this.env,
		"service": this.service,
		"address": this.address,
	}, sm)

	end := time.Now()
	if end.Unix()-start.Unix() > 60000 {
		log.Info(ctx, "show entering time for logging", log.Fields{
			"start": start,
			"end":   end,
			"order": 2,
		})
	}
	return resp, err
}

func initAccessLog(req interface{}, info *grpc.UnaryServerInfo) *log.AccessLog {
	accessLog := log.NewAccessLog()
	bytes, _ := json.Marshal(req)
	accessLog.Method = "POST"
	accessLog.Url = info.FullMethod
	accessLog.Body = string(bytes)
	return accessLog
}

func recordAccessLog(ctx context.Context, accessLog *log.AccessLog, err error, conf map[string]string, extra *sync.Map) {
	accessLog.RemoteAddress = util.ClientIP(ctx)
	accessLog.RemotePort = util.ClientPort(ctx)
	accessLog.StatusCode = getStatusCode(grpc.Code(err))
	accessLog.RequestId = util.ExtractRequestIDFromCtx(ctx)
	accessLog.TenantId = util.GetAccountId(ctx)
	accessLog.Others["category"] = "blogrpc"
	accessLog.Others["env"] = conf["env"]
	accessLog.Others["service"] = conf["service"]
	accessLog.Others["address"] = conf["address"]
	accessLog.Others["goroutines"] = runtime.NumGoroutine()
	if err != nil {
		accessLog.Others["error"] = err.Error()
	}
	if extra != nil {
		result := map[string]interface{}{}
		extra.Range(func(k, v interface{}) bool {
			key, ok := k.(string)
			if ok {
				result[key] = v
			}
			return true
		})
		if len(result) > 0 {
			accessLog.Others["extra"] = result
		}
	}

	accessLog.End()

	log.Stdout.Printf("%s", accessLog.Json())
}

func getStatusCode(code codes.Code) int {
	if code == codes.OK {
		return 200
	}

	// service customized error
	if code == codes.Unknown {
		return 400
	}

	//rpc internal error
	return 500
}
