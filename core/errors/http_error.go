package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strings"

	core_util "blogrpc/core/util"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type errorBody struct {
	Code      int32                  `protobuf:"varint,1,name=code" json:"code,omitempty"`
	Message   string                 `protobuf:"bytes,2,name=message" json:"message"`
	Errors    map[string]interface{} `protobuf:"bytes,3,name=errors" json:"errors,omitempty"`
	RequestId string                 `protobuf:"bytes,4,name=requestId" json:"requestId,omitempty"`
}

// Make this also conform to proto.Message for builtin JSONPb Marshaler
func (e *errorBody) Reset()         { *e = errorBody{} }
func (e *errorBody) String() string { return proto.CompactTextString(e) }
func (*errorBody) ProtoMessage()    {}

func HTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"message": "failed to marshal error message"}`

	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	s, body := parseErr(err)
	body.RequestId = w.Header().Get(core_util.RequestIDKey)
	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Printf("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Printf("Failed to write response: %v", err)
		}
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Printf("Failed to extract ServerMetadata from context")
	}

	handleForwardResponseServerMetadata(w, md)
	handleForwardResponseTrailerHeader(w, md)
	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Printf("Failed to write response: %v", err)
	}

	handleForwardResponseTrailer(w, md)
}

// CustomHTTPError Custom HTTP error handler
func CustomHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	s, body := parseErr(err)
	body.RequestId = w.Header().Get(core_util.RequestIDKey)
	w.Header().Set("Content-Type", marshaler.ContentType())

	// 分割字符串
	parts := strings.Split(body.Message, "\n\n")
	if len(parts) < 1 {
		fmt.Println("Invalid input format")
		return
	}

	// 提取各部分
	var errorCode, errorMsg, errorDetailsJSON string
	for i := range parts {
		if i == 0 {
			errorCode = parts[i]
		}
		if i == 1 {
			errorMsg = parts[i]
		}
		if i == 2 {
			errorDetailsJSON = parts[i]
		}
	}

	if errorCode == "1000001" || errorCode == "1000004" {
		w.WriteHeader(runtime.HTTPStatusFromCode(codes.Internal))
	} else if s.Code() != codes.Unknown {
		w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	} else {
		w.WriteHeader(runtime.HTTPStatusFromCode(codes.InvalidArgument))
	}

	customError := map[string]interface{}{
		"code":      errorCode,
		"message":   errorMsg,
		"errors":    errorDetailsJSON,
		"requestId": w.Header().Get(core_util.RequestIDKey),
	}
	marshaler.NewEncoder(w).Encode(customError)
}

func newInvalidArgumentStatus() *status.Status {
	return status.New(codes.InvalidArgument, "")
}

func newUnknowStatus() *status.Status {
	return status.New(codes.Unknown, "")
}

func parseErr(err error) (*status.Status, *errorBody) {
	if s, ok := status.FromError(err); ok {
		return s, &errorBody{Message: s.Message()}
	}

	if rpcErr, ok := err.(*RPCError); ok {
		return newInvalidArgumentStatus(), &errorBody{
			Message: rpcErr.Desc,
			Code:    int32(rpcErr.Code),
			Errors:  rpcErr.Extra,
		}
	}

	return newUnknowStatus(), &errorBody{Message: err.Error()}
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vs := range md.HeaderMD {
		if h, ok := outgoingHeaderMatcher(k); ok {
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
	}
}

func handleForwardResponseTrailerHeader(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k := range md.TrailerMD {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tKey)
	}
}

func handleForwardResponseTrailer(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}

func outgoingHeaderMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}

func OtherErrorHandler(w http.ResponseWriter, _ *http.Request, msg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	body, _ := json.Marshal(&errorBody{Message: msg})
	w.Write(body)
}
