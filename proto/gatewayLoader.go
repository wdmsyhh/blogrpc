package proto

import (
	"blogrpc/proto/hello"
	"blogrpc/proto/member"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
)

func NewGateway(ctx context.Context) (http.Handler, error) {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseEnumNumbers: true,
			UseProtoNames:  true,
		},
	}), runtime.WithIncomingHeaderMatcher(headerMatcherFunc))

	var err error

	err = hello.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, ":8081", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = member.RegisterMemberServiceHandlerFromEndpoint(ctx, mux, ":8082", []grpc.DialOption{grpc.WithInsecure()})
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
