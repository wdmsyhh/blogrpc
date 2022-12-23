package log

import (
	"os"

	"blogrpc/core/util"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	traceKey     = "backtrace"
	contextKey   = "context"
	categoryKey  = "category"
	requestIdKey = "reqId"
	TenantIdKey  = "tenantId"
	MaxStackSize = 4096
)

type Fields map[string]interface{}

var defaultFields logrus.Fields = logrus.Fields{"type": "service"}

// InitLogger initialize the default configuration for logger and change behavior based on envrionment
func InitLogger(level, env, service string) {
	l, err := logrus.ParseLevel(level)
	if nil != err {
		logrus.Fatalf("Failed to parse logger level setting, %v", err)
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetLevel(l)

	if env == "local" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, DisableColors: false, FullTimestamp: true})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				//change logrus default message key 'msg' to 'message'
				logrus.FieldKeyMsg: "message",
			},
		})
		logrus.SetOutput(os.Stdout)
		defaultFields[categoryKey] = service
	}
}

func Info(ctx context.Context, msg string, extra Fields) {
	newTraceEntry(ctx, extra).Info(msg)
}

func Warn(ctx context.Context, msg string, extra Fields) {
	newTraceEntry(ctx, extra).Warn(msg)
}

func WarnTrace(ctx context.Context, msg string, extra Fields, trace []byte) {
	newTraceEntry(ctx, extra).WithField(traceKey, string(trace)).Warn(msg)
}

func Error(ctx context.Context, msg string, extra Fields) {
	newTraceEntry(ctx, extra).Error(msg)
}

func ErrorTrace(ctx context.Context, msg string, extra Fields, trace []byte) {
	newTraceEntry(ctx, extra).WithField(traceKey, string(trace)).Error(msg)
}

func Panic(ctx context.Context, msg string, extra Fields, err interface{}) {
	newTraceEntry(ctx, extra).Error(msg)
	panic(err)
}

func newTraceEntry(ctx context.Context, extra Fields) *logrus.Entry {
	entry := logrus.WithFields(defaultFields)
	return entry.WithFields(logrus.Fields{
		contextKey:   extra,
		requestIdKey: util.ExtractRequestIDFromCtx(ctx),
		TenantIdKey:  util.GetAccountId(ctx),
	})
}
