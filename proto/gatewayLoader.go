package proto

import (
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
)

func NewGateway(ctx context.Context) (http.Handler, error) {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		OrigName:     true,
		EmitDefaults: true,
	}), runtime.WithIncomingHeaderMatcher(headerMatcherFunc))

	var err error

	err = hello.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, "localhost:8081", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = member.RegisterMemberServiceHandlerFromEndpoint(ctx, mux, "localhost:8082", []grpc.DialOption{grpc.WithInsecure()})
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

	return "", false
}
