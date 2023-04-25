package client

import (
	"blogrpc/core/component"
	"blogrpc/core/errors"
	"blogrpc/core/interceptor"
	"blogrpc/core/log"
	"fmt"

	"google.golang.org/grpc/resolver"

	"reflect"
	"strings"
	"unsafe"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientConn struct {
	Conn *grpc.ClientConn
}

type GetClientFunc func(*ClientConn) interface{}

const (
	ACCOUNT_ID   = "aid"
	X_REQUEST_ID = "x-request-id"

	TIMEOUT_ERROR             = "The service takes too long to execute"
	ISTIO_PROXY_TIMEOUT_ERROR = "downstream duration timeout"
)

var (
	serviceConns                   = component.NewSafeMap()
	serviceWithCtxInterceptorConns = component.NewSafeMap()
	// RPCProxy is used to call RPC server on behalf of the client.
	RPCProxy RPCCallProxy = &callProxy{}
)

// RPCCallProxy is a proxy to communcate with the remote RPC server.
type RPCCallProxy interface {
	Call(GetClientFunc, string, context.Context, interface{}) (interface{}, *errors.RPCError)
}

type callProxy struct {
}

// Call makes request to the remote RPC server. It uses a hystrix command to invoke remote service.
// If the remote service is not healthy, an internal error will be returned.
func (proxy *callProxy) Call(getClientFunc GetClientFunc, servMethodName string, ctx context.Context, req interface{}) (interface{}, *errors.RPCError) {
	return CallRPC(getClientFunc, servMethodName, ctx, req)
}

func Run(servMethodName string, ctx context.Context, req interface{}) (interface{}, *errors.RPCError) {
	getClientFunc := GetClientByFuncName(servMethodName)
	return RPCProxy.Call(getClientFunc, servMethodName, ctx, req)
}

// ListServiceConns list all conns
func ListServiceConns() map[string]interface{} {
	return serviceConns.List()
}

// CloseServiceConns close all established conns
func CloseServiceConns() {
	for _, conn := range serviceConns.List() {
		conn.(*grpc.ClientConn).Close()
	}
}

// CallRPC is an easy way to call RPC service
func CallRPC(getClientFunc GetClientFunc, servMethodName string, ctx context.Context, req interface{}) (resp interface{}, err *errors.RPCError) {
	if getClientFunc == nil {
		err = errors.NewInternal(fmt.Sprintf("nil getClientFunc for servMethodName %s", servMethodName))
		return
	}
	// send grpc request should use NewOutgoingContext
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// use old context ctx may cause context canceled
		ctx = metadata.NewOutgoingContext(context.Background(), md)
	}

	defer func() {
		if x := recover(); x != nil {
			err = errors.ConvertRecoveryError(x)
			log.Error(ctx, err.Error(), log.Fields{
				"code":  err.Code,
				"desc":  err.Desc,
				"extra": err.Extra,
			})
		}
	}()

	//servMethodName should follow Service.Method format
	parts := strings.Split(servMethodName, ".")
	var serviceName, subServiceName, methodName, serviceIdentifier string
	switch len(parts) {
	case 2: // for methods like MemberService.GetMember
		serviceName = parts[0]
		methodName = parts[1]
		serviceIdentifier = serviceName
	case 3: // for methods like EcService.ProductService.GetProduct
		serviceName = parts[0]
		subServiceName = parts[1]
		methodName = parts[2]
		serviceIdentifier = fmt.Sprintf("%s.%s", serviceName, subServiceName)
	default:
		return nil, errors.NewInternal("Service method name format is not valid")
	}

	conn := getServiceConn(serviceName)
	rpcClientImp := getClientFunc(&ClientConn{Conn: conn})

	var errRPC error
	resp, errRPC = makeRPCCall(rpcClientImp, serviceIdentifier, methodName, ctx, req)
	if errRPC != nil {
		err = errors.ToRPCError(errRPC)
	}
	return
}

func makeRPCCall(clientImp interface{}, service, method string, ctx context.Context, req interface{}) (ret interface{}, err error) {
	callRPCFunc := func(c context.Context, r interface{}) (response interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.ConvertRecoveryError(r)
			}
		}()

		// make RPC call based on reflection
		vClient := reflect.ValueOf(clientImp)
		f := vClient.Elem().MethodByName(method)
		if !f.IsValid() {
			err = errors.NewInternal(fmt.Sprintf("cannot get method %s in service %s", method, service))
			return nil, err
		}
		resp := f.Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(r),
		})

		if !resp[1].IsNil() {
			err := resp[1].Interface().(error)
			log.Warn(c, fmt.Sprintf("Fail to call %s.%s", service, method), log.Fields{
				"request":  r,
				"response": resp[0].Interface(),
				"error":    err.Error(),
			})
			return nil, err
		}
		return resp[0].Interface(), nil
	}

	return callRPCFunc(ctx, req)
}

func getServiceConn(serviceName string) *grpc.ClientConn {
	conn := serviceConns.Get(serviceName)
	if conn == nil {
		conn = dialWithOpts(serviceName, grpc.WithInsecure())
		serviceConns.Set(serviceName, conn)
	}
	return conn.(*grpc.ClientConn)
}

func GetServiceConnAddress(serviceName string) (string, error) {
	addrConnes := getStructPtrUnExportedField(reflect.ValueOf(getServiceConn(serviceName)), "conns").MapKeys()
	if len(addrConnes) > 0 {
		return getStructPtrUnExportedField(addrConnes[0], "curAddr").Interface().(resolver.Address).Addr, nil
	}
	return "", errors.NewNotExistsError("conns")
}

func getStructPtrUnExportedField(source reflect.Value, fieldName string) reflect.Value {
	// 获取非导出字段反射对象
	v := source.Elem().FieldByName(fieldName)
	// 构建指向该字段的可寻址（addressable）反射对象
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

func dialWithOpts(serviceName string, opts ...grpc.DialOption) *grpc.ClientConn {
	target := fmt.Sprintf("blogrpc-%s:1701", strings.ToLower(strings.Replace(serviceName, "Service", "", -1)))
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		panic(errors.NewInternal(fmt.Sprintf("%s service is unreachable with %v", serviceName, err)))
	}
	return conn
}

func GetServiceConnWithCtxInterceptor(serviceName string) *grpc.ClientConn {
	conn := serviceWithCtxInterceptorConns.Get(serviceName)
	if conn == nil {
		conn = dialWithOpts(serviceName, grpc.WithUnaryInterceptor(interceptor.ClientCtxInterceptor), grpc.WithInsecure())
		serviceWithCtxInterceptorConns.Set(serviceName, conn)
	}
	return conn.(*grpc.ClientConn)
}
