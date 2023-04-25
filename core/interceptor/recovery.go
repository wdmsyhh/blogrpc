package interceptor

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"blogrpc/core/errors"
	"blogrpc/core/log"

	"github.com/spf13/cast"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func init() {
	AddInterceptor(&RecoveryInterceptor{})
}

type RecoveryInterceptor struct {
	env     string
	service string
	address string
}

func (this *RecoveryInterceptor) Name() string { return "Recovery" }
func (this *RecoveryInterceptor) InitWithConf(conf map[string]interface{}, debug bool) error {
	this.service = cast.ToString(conf["service"])
	this.address = cast.ToString(conf["addr"])
	this.env = cast.ToString(conf["env"])
	return nil
}

// Recovery interceptor to handle grpc panic
func (this *RecoveryInterceptor) Handle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			// log stack
			stack := make([]byte, log.MaxStackSize)
			stack = stack[:runtime.Stack(stack, false)]
			// if panic, set custom error to 'err', in order that client and sense it.
			err = errors.ConvertRecoveryError(r)
			accessLog := initAccessLog(req, info)
			recordAccessLog(ctx, accessLog, err, map[string]string{
				"env":     this.env,
				"service": this.service,
				"address": this.address,
			}, nil)
			if strings.Contains(err.Error(), "no configuration of hosts") {
				log.Warn(ctx, err.Error(), log.Fields{})
				return
			}
			log.ErrorTrace(ctx, fmt.Sprintf("Fail to invoke method %s", info.FullMethod), log.Fields{
				"error": err.Error(),
			}, stack)
		}
	}()

	resp, err = handler(ctx, req)

	end := time.Now()
	if end.Unix()-start.Unix() > 60000 {
		log.Info(ctx, "show entering time for recovery", log.Fields{
			"start": start,
			"end":   end,
		})
	}
	return resp, err
}
