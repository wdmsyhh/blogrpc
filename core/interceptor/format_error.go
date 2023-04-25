package interceptor

import (
	"blogrpc/core/codes"
	"blogrpc/core/errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func init() {
	AddInterceptor(&FormatErrorInterceptor{})
}

type FormatErrorInterceptor struct {
}

// Name returns the name of FormatErrorInterceptor.
func (*FormatErrorInterceptor) Name() string { return "FormatError" }

// InitWithConf initialize the FormatErrorInterceptor before it is used.
func (v *FormatErrorInterceptor) InitWithConf(conf map[string]interface{}, debug bool) error {

	return nil
}

func (v *FormatErrorInterceptor) Handle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		if _, ok := err.(*errors.RPCError); !ok {
			err = &errors.RPCError{
				Code: codes.UnknownError,
				Desc: err.Error(),
			}
		}
	}

	return
}
