package interceptor

import (
	"blogrpc/core/util"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func init() {
	// AddInterceptor(&AuthInterceptor{})
}

type AuthInterceptor struct{}

func (this *AuthInterceptor) Name() string                                               { return "Auth" }
func (this *AuthInterceptor) InitWithConf(conf map[string]interface{}, debug bool) error { return nil }

// Auth interceptor for grpc
func (this *AuthInterceptor) Handle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	accountId := util.GetAccountId(ctx)

	if accountId == "" {
		err = grpc.Errorf(codes.Internal, "invalid account ID")
		return nil, err
	}

	return handler(ctx, req)
}
