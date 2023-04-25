package errors

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json"

	"blogrpc/core/codes"

	"google.golang.org/grpc"
	grpc_codes "google.golang.org/grpc/codes"
)

// RPCError defines the status from an RPC.
type RPCError struct {
	Code  codes.Code
	Desc  string
	Extra map[string]interface{}
}

type ErrorField struct {
	Code    string `protobuf:"varint,1,name=code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,name=message" json:"message"`
}

// Error returns the error information
func (self RPCError) Error() string {
	errMsg := fmt.Sprintf("%d%s%s", self.Code, codes.SEPARATOR, self.Desc)
	if self.Extra != nil {
		buf, err := json.Marshal(self.Extra)
		if err == nil {
			errMsg = fmt.Sprintf("%s%s%s", errMsg, codes.SEPARATOR, string(buf))
		}
	}

	return errMsg
}

// NewRPCError returns system pre-defined error
func NewRPCError(code codes.Code, desc string) *RPCError {
	RPCError := RPCError{
		Code: code,
		Desc: desc,
	}

	return &RPCError
}

func NewRPCErrorWithExtra(code codes.Code, desc string, extra map[string]interface{}) *RPCError {
	RPCError := RPCError{
		Code:  code,
		Desc:  desc,
		Extra: extra,
	}

	return &RPCError
}

func ToRPCError(err error) *RPCError {
	if RPCError, ok := err.(*RPCError); ok {
		return RPCError
	}

	gRPCErrorCode := grpc.Code(err)
	gRPCErrorDesc := grpc.ErrorDesc(err)
	if gRPCErrorCode != grpc_codes.Unknown {
		return NewInternal(gRPCErrorDesc)
	}

	BlogrpcCodeStr, BlogrpcDesc, BlogrpcExtraStr := parseGrpcDesc(gRPCErrorDesc)
	BlogrpcCode, _ := strconv.ParseUint(BlogrpcCodeStr, 10, 0)
	var BlogrpcExtra map[string]interface{}
	json.Unmarshal([]byte(BlogrpcExtraStr), &BlogrpcExtra)
	return &RPCError{
		Code:  codes.Code(uint32(BlogrpcCode)),
		Desc:  BlogrpcDesc,
		Extra: BlogrpcExtra,
	}
}

func ConvertRecoveryError(err interface{}) *RPCError {
	BlogRPCError, ok := err.(*RPCError)
	if ok {
		return BlogRPCError
	} else {
		RPCError, ok := err.(error)
		if ok {
			return NewInternal(RPCError.Error())
		} else {
			return NewUnknowError(err)
		}
	}
}

func NewInternal(desc string) *RPCError {
	return &RPCError{
		Code: codes.InternalServerError,
		Desc: desc,
	}
}

func IsRPCInternalError(err error) bool {
	BlogrpcErr, ok := err.(*RPCError)
	return ok && BlogrpcErr.Code == codes.InternalServerError
}

func NewUnknowError(err interface{}) *RPCError {
	return &RPCError{
		Code: codes.UnknownError,
		Desc: fmt.Sprintf("unknown error: %v", err),
	}
}

func parseGrpcDesc(errorMsg string) (string, string, string) {
	codeEndIndex := strings.Index(errorMsg, codes.SEPARATOR)
	if codeEndIndex <= 0 {
		return "", "", ""
	}

	descStartIndex := codeEndIndex + len(codes.SEPARATOR)

	index := strings.Index(errorMsg[descStartIndex:len(errorMsg)], codes.SEPARATOR)
	if index <= 0 {
		return errorMsg[0:codeEndIndex], errorMsg[descStartIndex:len(errorMsg)], ""
	}

	RpcErrMsg := strings.Split(errorMsg, codes.SEPARATOR)

	return RpcErrMsg[0], RpcErrMsg[1], RpcErrMsg[2]
}

func NewInvalidParamsError(extraMsg map[string]interface{}) error {
	return NewRPCErrorWithExtra(codes.InvalidParams, "Invalid Params", extraMsg)
}

func NewNotExistsError(field string) error {
	return NewNotExistsErrorWithMessage(field, "")
}

func NewNotExistsErrorWithMessage(field string, message string) error {
	if message == "" {
		message = "Resource not exists"
	}
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "notExists",
			Message: message,
		},
	})
}

func NewNotEnabledError(field string) error {
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "notEnabled",
			Message: "Resource is not enabled yet",
		},
	})
}

func NewAlreadyExistsError(field string) error {
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "alreadyExists",
			Message: "Resource already exists",
		},
	})
}

func NewCanNotDeleteNonEmptyItemError(field string) error {
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "canNotDeleteNonEmptyItem",
			Message: "Resource is not empty",
		},
	})
}

func NewCanNotEditError(field string) error {
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "canNotEdit",
			Message: "Resource can not be edited",
		},
	})
}

func NewCanNotDeleteReferencedItemError(field string) error {
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "canNotDeleteReferencedItem",
			Message: "Resource is referenced",
		},
	})
}

func NewCanNotDeleteDefaultItemError(field string) error {
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "canNotDeleteDefaultItem",
			Message: "Resource is default",
		},
	})
}

func NewInvalidArgumentErrorWithMessage(field string, message string) error {
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "invalidArgument",
			Message: message,
		},
	})
}

func NewInvalidArgumentError(field string) error {
	return NewInvalidArgumentErrorWithMessage(field, "")
}

func NewTooManyRequestsError(field string) error {
	return NewInvalidParamsError(map[string]interface{}{
		field: ErrorField{
			Code:    "tooManyRequests",
			Message: "Too many requests",
		},
	})
}

func IsNotExistsError(err error) bool {
	return isError(err, "", "notExists")
}

func IsNotEnabledError(err error, field string) bool {
	return isError(err, field, "notEnabled")
}

func IsAlreadyExistsError(err error) bool {
	return isError(err, "", "alreadyExists")
}

func IsInvalidArgumentError(err error) bool {
	return isError(err, "", "invalidArgument")
}

func isError(err error, field, code string) bool {
	maiErr := ToRPCError(err)
	if maiErr == nil {
		return false
	}

	for f, data := range maiErr.Extra {
		// 当参数 field 不为空时，需要判断指定 field 的信息
		if field != "" && f != field {
			continue
		}
		b, e := json.Marshal(data)
		if e != nil {
			continue
		}
		errorField := ErrorField{}
		e = json.Unmarshal(b, &errorField)
		if e != nil {
			continue
		}

		if errorField.Code == code {
			return true
		}
	}

	return false
}
