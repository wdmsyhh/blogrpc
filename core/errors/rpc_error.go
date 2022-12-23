package errors

// TODO: Benny Zhou 12/01/2016 - file rpc_error.go should be renamed to Blogrpc_go
// after the existing file "Blogrpc_error.go" removed.

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json"

	"blogrpc/core/codes"

	"google.golang.org/grpc"
	grpc_codes "google.golang.org/grpc/codes"
)

// BlogrpcError defines the status from an RPC.
type BlogrpcError struct {
	Code  codes.Code
	Desc  string
	Extra map[string]interface{}
}

type ErrorField struct {
	Code    string `protobuf:"varint,1,name=code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,name=message" json:"message"`
}

// Error returns the error information
func (self BlogrpcError) Error() string {
	errMsg := fmt.Sprintf("%d%s%s", self.Code, codes.SEPARATOR, self.Desc)
	if self.Extra != nil {
		buf, err := json.Marshal(self.Extra)
		if err == nil {
			errMsg = fmt.Sprintf("%s%s%s", errMsg, codes.SEPARATOR, string(buf))
		}
	}

	return errMsg
}

// NewBlogrpcError returns system pre-defined error
func NewBlogrpcError(code codes.Code, desc string) *BlogrpcError {
	rpcError := BlogrpcError{
		Code: code,
		Desc: desc,
	}

	return &rpcError
}

func NewBlogrpcErrorWithExtra(code codes.Code, desc string, extra map[string]interface{}) *BlogrpcError {
	rpcError := BlogrpcError{
		Code:  code,
		Desc:  desc,
		Extra: extra,
	}

	return &rpcError
}

func ToblogrpcError(err error) *BlogrpcError {
	if BlogrpcError, ok := err.(*BlogrpcError); ok {
		return BlogrpcError
	}

	grpcErrorCode := grpc.Code(err)
	grpcErrorDesc := grpc.ErrorDesc(err)
	if grpcErrorCode != grpc_codes.Unknown {
		return NewInternal(grpcErrorDesc)
	}

	BlogrpcCodeStr, BlogrpcDesc, BlogrpcExtraStr := parseGrpcDesc(grpcErrorDesc)
	BlogrpcCode, _ := strconv.ParseUint(BlogrpcCodeStr, 10, 0)
	var BlogrpcExtra map[string]interface{}
	json.Unmarshal([]byte(BlogrpcExtraStr), &BlogrpcExtra)
	return &BlogrpcError{
		Code:  codes.Code(uint32(BlogrpcCode)),
		Desc:  BlogrpcDesc,
		Extra: BlogrpcExtra,
	}
}

func ConvertRecoveryError(err interface{}) *BlogrpcError {
	BlogrpcError, ok := err.(*BlogrpcError)
	if ok {
		return BlogrpcError
	} else {
		rpcError, ok := err.(error)
		if ok {
			return NewInternal(rpcError.Error())
		} else {
			return NewUnknowError(err)
		}
	}
}

func NewInternal(desc string) *BlogrpcError {
	return &BlogrpcError{
		Code: codes.InternalServerError,
		Desc: desc,
	}
}

func IsRPCInternalError(err error) bool {
	BlogrpcErr, ok := err.(*BlogrpcError)
	return ok && BlogrpcErr.Code == codes.InternalServerError
}

func NewUnknowError(err interface{}) *BlogrpcError {
	return &BlogrpcError{
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

	BlogrpcErrMsg := strings.Split(errorMsg, codes.SEPARATOR)

	return BlogrpcErrMsg[0], BlogrpcErrMsg[1], BlogrpcErrMsg[2]
}

func NewInvalidParamsError(extraMsg map[string]interface{}) error {
	return NewBlogrpcErrorWithExtra(codes.InvalidParams, "Invalid Params", extraMsg)
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
	maiErr := ToblogrpcError(err)
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
