package controller

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"blogrpc/core/log"
	rpc_util "blogrpc/core/util"
	"blogrpc/openapi/business/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/serenize/snaker"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	DEFAULT_PAGE              = 1
	DEFAULT_PER_PAGE          = 10
	StatusUnprocessableEntity = 422
	MODULE_API_PATH           = "/modules/*rest"
	APP_API_PATH              = "/apps/:appId/*rest"
	DEBUG_PPROF_PATH          = "/debug/pprof/*rest"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.SetTagName("valid")
	validate.RegisterAlias("Required", "required")
}

type IPager interface {
	GetCurPage() int32
	GetPerPage() int32

	SetCurPage(page int32)
	SetPerPage(perPage int32)
}

// Common Request Payload
type FormPager struct {
	Page    uint32 `query:"page"`
	PerPage uint32 `query:"per_page"`
}

type FormSort struct {
	Sort      string `query:"sort"`
	Direction string `query:"direction"`
}

func (self *FormPager) GetCurPage() uint32 {
	if 0 < self.Page {
		return self.Page
	}

	return DEFAULT_PAGE //default from page 1
}

func (self *FormPager) GetPerPage() uint32 {
	if 0 < self.PerPage && self.PerPage <= 100 {
		return self.PerPage
	}

	return DEFAULT_PER_PAGE //default returns 10 items
}

func (self *FormSort) GetSort(sortMap map[string]string) string {
	var sort string
	if nil != sortMap && "" != sortMap[self.Sort] {
		sort = sortMap[self.Sort]
	} else {
		sort = "createdAt"
	}
	if strings.EqualFold(self.Direction, "asc") {
		sort = "+" + sort
	} else {
		sort = "-" + sort
	}
	return sort
}

func (self *FormPager) SetCurPage(page uint32) {
	self.Page = page
}

func (self *FormPager) SetPerPage(page uint32) {
	self.PerPage = page
}

type RespList struct {
	Items   interface{} `json:"items"`
	Page    uint32      `json:"page,omitempty"`
	PerPage uint32      `json:"per_page,omitempty"`
	Total   uint64      `json:"total"`
}

type RespListWithoutPagination struct {
	Items interface{} `json:"items"`
	Total uint64      `json:"total"`
}

type MyHandlerFunc func(*gin.Context)

type ControllerAction struct {
	Name        string
	Method      string
	Path        string
	Scope       string // Add scope to validate the caller permission, if empty means public access
	HandlerFunc MyHandlerFunc
	IsPublic    bool
	FixPath     bool
}

func (action *ControllerAction) GenerateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer reco(c)

		for _, pathParam := range c.Params {
			if "" == pathParam.Value {
				c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": fmt.Sprintf("missing path param %s", pathParam.Key)})
				return
			}
		}

		action.HandlerFunc(c)
	}
}

func reco(c *gin.Context) {
	if err := recover(); err != nil {
		baseErr, ok := err.(util.BaseError)
		if ok {
			baseErr.Handle(c)
		} else {
			panic(err) //panic if the error is not predefined
		}
	}
}

var (
	actions = make([]*ControllerAction, 0, 100)
)

func GetAllControllerAction() []*ControllerAction {
	return actions
}

func enableActions(c ...*ControllerAction) {
	if c != nil {
		actions = append(actions, c...)
	} else {
		log.Error(nil, "Failed to enable actions, action is nil.", nil)
	}
}

func panicWithInvalidParamError() {
	panic(util.NewInvalidParamError("10004", "Param missing"))
}

func GetGrpcContext(c *gin.Context) context.Context {
	//Refer: https://godoc.org/google.golang.org/grpc/metadata
	ctx := context.Background()

	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.New(map[string]string{
			rpc_util.AccountIdKey:     getAccountId(c),
			rpc_util.RequestIDKey:     c.MustGet(rpc_util.RequestIDKey).(string),
			rpc_util.TracingHeaderKey: c.Request.Header.Get(rpc_util.TracingHeaderKey),
		}),
	)

	return ctx
}

func writeJson(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// parsePayload: parse content form URL query or request body
// interface obj must be a pointer to struct
func parsePayload(c *gin.Context, obj interface{}) {
	req := c.Request

	// Parse query string
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("query")
	decoder.IgnoreUnknownKeys(true)
	q := req.URL.Query()
	if err := decoder.Decode(obj, q); nil != err {
		log.Warn(c, "Failed to decode query string", log.Fields{
			"url":    req.RequestURI,
			"method": req.Method,
			"query":  q,
			"error":  err,
		})
		util.PanicInvalidBinding(http.StatusBadRequest, map[string]string{"message": "Problems parsing query string values."})
	}

	handleDecodeJsonBody(c, obj)
}

// handleDecodeJsonBody call decodeJsonBody func and handle error when decode json body
func handleDecodeJsonBody(c *gin.Context, obj interface{}) {
	req := c.Request

	//only parse body to JSON when Conent-Length > 0
	if req.ContentLength > 0 {
		body, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if err := DecodeJsonBody(req, obj); err != nil {
			dumpReq, _ := httputil.DumpRequest(req, true)

			log.Warn(c, "Received invalid JSON", log.Fields{
				"url":     req.RequestURI,
				"method":  req.Method,
				"request": string(dumpReq),
				"error":   err,
				"body":    string(body),
			})
			util.PanicInvalidBinding(http.StatusBadRequest, map[string]string{"message": "Problems parsing JSON"})
		}
	}
}

// decodeJsonBody reads the request body and decodes the JSON using json.NewDecoder.
func DecodeJsonBody(request *http.Request, v interface{}) error {
	decoder := json.NewDecoder(request.Body)
	return decoder.Decode(v)
}

func DecodeFormBody(c *gin.Context, v interface{}) error {
	req := c.Request
	err := req.ParseForm()
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	for k := range req.PostForm {
		val, ok := c.GetPostFormArray(k)
		if ok {
			switch len(val) {
			case 0:
				// empty
			case 1:
				data[k] = val[0]
			default:
				data[k] = val
			}
		}
	}
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	json.Unmarshal(result, v)
	return nil
}

func validatePayload(obj interface{}) {
	err := validate.Struct(obj)
	if err == nil {
		return
	}

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return
	}

	code := "missing_field"
	errMsg := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := snaker.CamelToSnake(err.Field())
		if "required" == strings.ToLower(err.Tag()) {
			errMsg[field] = "Can not be empty"
		} else {
			errMsg[field] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", field, err.Tag())
			code = "invalid"
		}
	}
	errMsg["code"] = code
	util.PanicInvalidBinding(StatusUnprocessableEntity, map[string]interface{}{"message": "Validation Failed", "errors": errMsg})
}

func getAccountId(c *gin.Context) string {
	return c.MustGet(rpc_util.AccountIdKey).(string)
}

func BindPayload(c *gin.Context, payload interface{}) {
	parsePayload(c, payload)

	if pager, ok := payload.(IPager); ok {
		pager.SetCurPage(pager.GetCurPage()) //initialize the default value
		pager.SetPerPage(pager.GetPerPage()) //initialize the default value
		//above 2 lines of code, make sure the page and cur_page can pass the validation rule
	}

	validatePayload(payload)
}

func BindMapAndPayload(c *gin.Context, m *map[string]interface{}, payload interface{}) {
	handleDecodeJsonBody(c, m)

	util.DecodeMapToStruct(*m, payload)

	validatePayload(payload)
}

func BindMapAndPayloadByWeaklyTyped(c *gin.Context, m *map[string]interface{}, payload interface{}, mHandler func(*map[string]interface{})) {
	handleDecodeJsonBody(c, m)

	if mHandler != nil {
		mHandler(m)
	}

	if err := util.DecodeMapToStructByWeaklyTyped(*m, payload); err != nil {
		log.Error(c, "failed to decode map to struct by weakly typed input", log.Fields{
			"url":    c.Request.RequestURI,
			"method": c.Request.Method,
			"error":  err.Error(),
			"m":      fmt.Sprintf("%#v", m),
		})
		util.PanicInvalidBinding(http.StatusBadRequest, map[string]string{"message": "Parsing Failed."})
	}

	validatePayload(payload)
}

func BindXmlAndPayload(c *gin.Context, payload interface{}) {
	req := c.Request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warn(c, "Received invalid xml", log.Fields{
			"url":    req.RequestURI,
			"method": req.Method,
			"error":  err,
			"body":   string(body),
		})
	}

	if err := xml.Unmarshal(body, payload); err != nil {
		util.PanicInvalidBinding(http.StatusBadRequest, map[string]string{"message": "Problems parsing XML"})
	}
}
