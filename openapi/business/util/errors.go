package util

import (
	"net/http"
	"runtime/debug"

	"blogrpc/core/log"

	"github.com/gin-gonic/gin"
)

type BaseError interface {
	Handle(c *gin.Context)
}

type ApiError struct {
	Code          int
	Msg           string
	InternalError error
}

func (err *ApiError) Handle(c *gin.Context) {
	log.ErrorTrace(c, err.Msg, log.Fields{
		"code":  err.Code,
		"error": err.InternalError,
	}, debug.Stack())
	HttpError(c, map[string]string{"message": err.Msg}, http.StatusInternalServerError)
}

func NewApiError(code int, msg string, err error) *ApiError {
	return &ApiError{Code: code, Msg: msg, InternalError: err}
}

const (
	ServiceError = 50000 + iota //5 means service package error
	ServiceJSONError
)

const (
	ControllerError = 60000 + iota //6 means controller package error
	ControllerFormPayloadError
)

const (
	ModelError = 70000 + iota
	ModelDBError
	RedisError
)

type WebhookMissingParam struct {
	Msg string
}

func (err *WebhookMissingParam) Handle(c *gin.Context) {
	log.ErrorTrace(c, err.Msg, log.Fields{
		"message": err.Msg,
	}, debug.Stack())
	c.JSON(http.StatusOK, map[string]string{"message": err.Msg})
}

func NewWarningWebhookMissingParam(msg string) *WebhookMissingParam {
	return &WebhookMissingParam{Msg: msg}
}

type InvalidParamError struct {
	Msg  string `json:"message"`
	Code string `json:"code,omitempty"`
}

func (err *InvalidParamError) Handle(c *gin.Context) {
	log.WarnTrace(c, err.Msg, log.Fields{
		"code": err.Code,
	}, debug.Stack())
	HttpError(c, err, 422)
}

func NewInvalidParamError(code string, msg string) *InvalidParamError {
	return &InvalidParamError{Code: code, Msg: msg}
}

type UnexpectedWarning struct {
	Resonse interface{}
}

func (this *UnexpectedWarning) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, this.Resonse)
}

func NewUnexpectedWarning(response interface{}) *UnexpectedWarning {
	return &UnexpectedWarning{
		Resonse: response,
	}
}

type InvalidPayloadBinding struct {
	statusCode int
	msg        interface{}
}

func (this *InvalidPayloadBinding) Handle(c *gin.Context) {
	c.JSON(this.statusCode, this.msg)
}

func PanicInvalidBinding(statusCode int, msg interface{}) {
	panic(&InvalidPayloadBinding{statusCode, msg})
}
