package proto

import (
	"context"
	"net/http"
	"time"

	"blogrpc/core/constant"
	"blogrpc/core/errors"
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func NewGateway(ctx context.Context) (http.Handler, error) {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		OrigName:     true,
		EmitDefaults: true,
	}), runtime.WithIncomingHeaderMatcher(headerMatcherFunc),
		runtime.WithProtoErrorHandler(errors.CustomHTTPError))

	//runtime.HTTPError = errors.HTTPError2
	//runtime.OtherErrorHandler = errors.OtherErrorHandler
	runtime.DefaultContextTimeout = 10 * time.Second

	var err error

	err = hello.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, constant.ServiceHelloHost+":"+constant.ServiceHelloPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = member.RegisterMemberServiceHandlerFromEndpoint(ctx, mux, constant.ServiceMemberHost+":"+constant.ServiceMemberPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	return mux, nil
}

func headerMatcherFunc(headerName string) (string, bool) {
	if headerName == "App-Id" {
		return headerName, true
	}

	if headerName == "App-Secret" {
		return headerName, true
	}

	if headerName == "Aid" {
		return headerName, true
	}

	return "", false
}
