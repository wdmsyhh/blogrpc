package interceptor

import (
	"blogrpc/core/errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func getOutCtx(incoming context.Context) context.Context {
	ctx := incoming
	if md, ok := metadata.FromIncomingContext(incoming); ok {
		// use old context ctx may cause context canceled
		ctx = metadata.NewOutgoingContext(context.Background(), md)
	}
	return ctx
}

func ClientCtxInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	outCtx := getOutCtx(ctx)
	errRPC := invoker(outCtx, method, req, reply, cc, opts...)
	if errRPC != nil {
		return errors.ToRPCError(errRPC)
	}
	return nil
}
