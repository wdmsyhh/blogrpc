package extension

import (
	"blogrpc/core/log"
	"blogrpc/core/util"
	sysContext "context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	goLog "log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cast"
	"golang.org/x/net/context"
)

const (
	DEFAULT_TIMEOUT      = 10               // seconds
	RETRY_COUNT          = 3                // total number of retries
	RETRY_TIME_INCR      = 10 * time.Second // time increment for each retry
	RETRY_SWITCH         = "RetrySwitch"    // retry switch key
	REQUEST_TIMEOUT      = "RequestTimeout" // request timeout key
	BOUNCE_TO_RAW_STRING = "BounceToRawString"
)

var (
	// RequestClient is the proxy instance for actaul request instance
	RequestClient Requestor
	restClient    *RestClient
	zeroDialer    net.Dialer
)

func init() {
	request := NewRestClient()
	RequestClient = request
	RegisterExtension(request)
}

type RestClient struct {
	name       string
	config     map[string]interface{}
	cookies    []*http.Cookie
	Debug      bool
	Timeout    int64
	Logger     *goLog.Logger
	agent      *gorequest.SuperAgent
	RemoteAddr string // record server's address after each http request
}

// Requestor sends request as client proxy to remote server, and every function that need to request remote server
// have a tracing version, which will be traced during the request lifetime. For those non-tracing version functions,
// use TODO context to aviod being traced.
type Requestor interface {
	Get(ctx context.Context, service string, path string, params *url.Values, pHeaders *map[string]string) ([]byte, gorequest.Response, error)
	PostJson(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error)
	PutJson(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error)
	PatchJson(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error)
	DeleteJson(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error)
	PostFormData(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error)
	PostFile(ctx context.Context, service, path, fieldName string, fileBuf []byte, pHeaders *map[string]string) ([]byte, gorequest.Response, error)
	RequestJson(ctx context.Context, method string, fullUrl string, data interface{}, pHeaders *map[string]string, params *url.Values) ([]byte, gorequest.Response, error)
	RequestXml(ctx context.Context, method string, fullUrl string, data interface{}, pHeaders *map[string]string, params *url.Values) ([]byte, gorequest.Response, error)
	GetFullUrl(service string, path string) string
	SetCookies(cookies []*http.Cookie)
	Close()
	Clone() Requestor
}

func NewRestClient() *RestClient {
	if restClient == nil {
		return &RestClient{
			config:  make(map[string]interface{}),
			name:    "request",
			cookies: make([]*http.Cookie, 0),
			Debug:   false,
			Timeout: DEFAULT_TIMEOUT,
			Logger:  nil,
		}
	}
	return restClient
}

func (self *RestClient) Name() string {
	return self.name
}

func (self *RestClient) Clone() Requestor {
	newClient := *self
	return &newClient
}

func (self *RestClient) InitWithConf(conf map[string]interface{}, debug bool) error {
	if timeout, ok := conf["timeout"]; ok {
		self.Timeout = timeout.(int64)
	}
	self.config = conf
	self.Debug = debug
	self.Logger = getLogger("[httpclient] ")
	self.agent = self.newHttpClient()
	return nil
}

func (self *RestClient) newHttpClient() *gorequest.SuperAgent {
	httpClient := gorequest.New()
	httpClient.SetDebug(self.Debug)
	httpClient.Client.Timeout = time.Duration(self.Timeout) * time.Second
	httpClient.Transport.DisableKeepAlives = false
	httpClient.Transport.MaxIdleConnsPerHost = 100
	httpClient.Transport.IdleConnTimeout = 30 * time.Second

	// rewrite DialContext to get server's address
	httpClient.Transport.DialContext = func(ctx sysContext.Context, network, addr string) (net.Conn, error) {
		conn, err := zeroDialer.DialContext(ctx, network, addr)
		if conn != nil {
			self.RemoteAddr = conn.RemoteAddr().String()
		}
		return conn, err
	}

	return httpClient
}

func (self *RestClient) copyHttpClient(ctx context.Context) *gorequest.SuperAgent {
	client := self.newHttpClient()
	client.SetLogger(self.Logger)
	if self.agent != nil {
		client.Transport = self.agent.Transport
		seconds := ctx.Value(REQUEST_TIMEOUT)
		if seconds == nil {
			client.Client = self.agent.Client
			return client
		}
		client.Client.Timeout = time.Second * time.Duration(seconds.(int64))
	}
	return client
}

func (self *RestClient) Close() {
}

func (self *RestClient) SetCookies(cookies []*http.Cookie) {
	self.cookies = cookies
}

func (self *RestClient) Get(ctx context.Context, service string, path string, params *url.Values, pHeaders *map[string]string) ([]byte, gorequest.Response, error) {
	return self.RequestJson(ctx, "GET", self.GetFullUrl(service, path), nil, pHeaders, params)
}

func (self *RestClient) PostJson(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error) {
	return self.RequestJson(ctx, "POST", self.GetFullUrl(service, path), data, pHeaders, nil)
}

func (self *RestClient) PutJson(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error) {
	return self.RequestJson(ctx, "PUT", self.GetFullUrl(service, path), data, pHeaders, nil)
}

func (self *RestClient) PatchJson(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error) {
	return self.RequestJson(ctx, "PATCH", self.GetFullUrl(service, path), data, pHeaders, nil)
}

func (self *RestClient) DeleteJson(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error) {
	return self.RequestJson(ctx, "DELETE", self.GetFullUrl(service, path), data, pHeaders, nil)
}

func (self *RestClient) PostFormData(ctx context.Context, service string, path string, data interface{}, pHeaders *map[string]string) ([]byte, gorequest.Response, error) {
	var err error
	fullUrl := path
	if domain, ok := self.config["domain-"+service]; ok {
		fullUrl = domain.(string) + path
	}

	log.Info(ctx, "Send POST request with PostFormData method", log.Fields{
		"service": service,
		"url":     fullUrl,
		"data":    data,
	})

	httpClient := self.copyHttpClient(ctx)
	sa := httpClient.Post(fullUrl).Type("urlencoded").Send(data)
	if pHeaders == nil {
		pHeaders = &map[string]string{
			"Accept":       "application/json;charset=utf-8",
			"Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
		}
	}
	//Support add headers
	self.setHeaders(ctx, httpClient, pHeaders)

	// used to record the duration of request
	// so need to initialize it before calling sa.End()
	startTime := time.Now()
	u, err := url.Parse(fullUrl)
	if err != nil {
		u = &url.URL{}
	}
	accessOutLog := log.NewAccessOutLog(
		util.ExtractRequestIDFromCtx(ctx),
		u.Host,
		u.Scheme,
		"POST",
		u.Path,
		cast.ToString(data),
		util.GetAccountId(ctx),
	)

	resp, body, errs := sa.End()
	statusCode := 599 // default code for request that have no response
	if resp != nil {
		statusCode = resp.StatusCode
	}
	accessOutLog.End(ctx, statusCode, self.RemoteAddr, body, startTime, collectErrors(errs))

	if errs != nil && len(errs) > 0 {
		return []byte{}, resp, collectErrors(errs)
	}

	bytes := ([]byte)(body)
	return bytes, resp, err
}

func (self *RestClient) PostFile(ctx context.Context, service, path, fieldName string, fileBuf []byte, pHeaders *map[string]string) ([]byte, gorequest.Response, error) {
	var err error
	fullUrl := path
	if domain, ok := self.config["domain-"+service]; ok {
		fullUrl = domain.(string) + path
	}

	log.Info(ctx, "Send POST request with PostFormData method", log.Fields{
		"service": service,
		"url":     fullUrl,
	})

	httpClient := self.copyHttpClient(ctx)
	sa := httpClient.Post(fullUrl).Type("multipart").SendFile(fileBuf, "", fieldName)

	//Support add headers
	self.setHeaders(ctx, httpClient, pHeaders)

	// used to record the duration of request
	// so need to initialize it before calling sa.End()
	startTime := time.Now()
	u, err := url.Parse(fullUrl)
	if err != nil {
		u = &url.URL{}
	}
	accessOutLog := log.NewAccessOutLog(
		util.ExtractRequestIDFromCtx(ctx),
		u.Host,
		u.Scheme,
		"POST",
		u.Path,
		"",
		util.GetAccountId(ctx),
	)

	resp, body, errs := sa.End()
	statusCode := 599 // default code for request that have no response
	if resp != nil {
		statusCode = resp.StatusCode
	}
	accessOutLog.End(ctx, statusCode, self.RemoteAddr, body, startTime, collectErrors(errs))

	if errs != nil && len(errs) > 0 {
		return []byte{}, resp, collectErrors(errs)
	}

	bytes := ([]byte)(body)
	return bytes, resp, err
}

func (self *RestClient) RequestJson(ctx context.Context, method string, fullUrl string, data interface{}, pHeaders *map[string]string, params *url.Values) ([]byte, gorequest.Response, error) {
	var paramStr string = ""
	var jsonData []byte
	var err error

	if method != "get" && method != "GET" {
		if v := reflect.ValueOf(data); v.Kind() == reflect.String {
			jsonData = []byte(v.String())
		} else {
			jsonData, err = json.Marshal(data)
			if nil != err {
				log.Panic(ctx, "Fail to convert struct to JSON", log.Fields{
					"data":     data,
					"error":    err,
					"errorMsg": err.Error(),
				}, err)
			}
		}
	}

	httpClient := self.copyHttpClient(ctx)
	switch method {
	case "GET", "get":
		httpClient.Get(fullUrl)
	case "POST", "post":
		httpClient.Post(fullUrl).Send(string(jsonData))
	case "DELETE", "delete":
		httpClient.Delete(fullUrl).Send(string(jsonData))
	case "PUT", "put":
		httpClient.Put(fullUrl).Send(string(jsonData))
	case "PATCH", "patch":
		httpClient.Patch(fullUrl).Send(string(jsonData))

	}

	if bounceToRawString := ctx.Value(BOUNCE_TO_RAW_STRING); bounceToRawString != nil {
		httpClient.BounceToRawString = bounceToRawString.(bool)
	}

	//Parse parameters as JSON
	if params != nil {
		paramStr = params.Encode()
		httpClient.Query(paramStr)
	}

	//Add cookies if has
	if len(self.cookies) > 0 {
		httpClient.AddCookies(self.cookies)
	}

	//Support add headers
	self.setHeaders(ctx, httpClient, pHeaders)

	// used to record the duration of request
	// so need to initialize it before calling httpClient.End()
	startTime := time.Now()
	u, err := url.Parse(fullUrl)
	if err != nil {
		u = &url.URL{}
	}
	urlPath := u.Path
	if params != nil {
		encoded := params.Encode()
		if encoded != "" {
			urlPath = urlPath + "?" + encoded
		}
	}
	// sometimes we'll add parameter in fullUrl
	urlParams := u.Query().Encode()
	if urlParams != "" {
		if strings.Contains(urlPath, "?") {
			urlPath = urlPath + "&" + urlParams
		} else {
			urlPath = urlPath + "?" + urlParams
		}
	}
	accessOutLog := log.NewAccessOutLog(
		util.ExtractRequestIDFromCtx(ctx),
		u.Host,
		u.Scheme,
		strings.ToUpper(method),
		urlPath,
		string(jsonData),
		util.GetAccountId(ctx),
	)

	resp, body, errs := httpClient.End()
	statusCode := 599 // default code for request that have no response
	if resp != nil {
		statusCode = resp.StatusCode
	}
	accessOutLog.End(ctx, statusCode, self.RemoteAddr, body, startTime, collectErrors(errs))

	if errs != nil {
		if isRetryable(ctx, errs) {
			errs = retry(RETRY_COUNT, RETRY_TIME_INCR, func() []error {
				startTime := time.Now()
				resp, body, errs = httpClient.End()
				statusCode := 599 // default code for request that have no response
				if resp != nil {
					statusCode = resp.StatusCode
				}
				accessOutLog.End(ctx, statusCode, self.RemoteAddr, body, startTime, collectErrors(errs))
				return errs
			}, isNetTimeout)
		}

		if errs != nil {
			return []byte{}, resp, collectErrors(errs)
		}
	}

	bytes := ([]byte)(body)
	return bytes, resp, err
}

// retryCallback is the business logic callback to be retried.
type retryCallback func() []error

// retryPredicate is a predicate assertion whether the retry fails or need to retry again.
type retryPredicate func(errs []error) bool

// retry is used to execute the retry logic according to the retry policy passed.
func retry(retryCount int, retryTime time.Duration, callback retryCallback, predicate retryPredicate) (errs []error) {
	for attempt := 0; attempt < retryCount; attempt++ {
		time.Sleep(time.Duration(attempt) * retryTime)

		errs = callback()
		if errs == nil || !predicate(errs) {
			return
		}
	}
	return []error{
		fmt.Errorf("after %d attempts, last errors: %s", retryCount, errs),
	}
}

// OpenRequestFailedRetrySwitch is used to open request failed retry switch in request context.
func OpenRequestFailedRetrySwitch(ctx context.Context) context.Context {
	return context.WithValue(ctx, RETRY_SWITCH, true)
}

func isRetryable(ctx context.Context, errs []error) bool {
	switchStatus := ctx.Value(RETRY_SWITCH)
	return switchStatus != nil && switchStatus.(bool) && isNetTimeout(errs)
}

func isNetTimeout(errs []error) bool {
	if err, ok := errs[0].(*url.Error); ok {
		if err, ok := err.Err.(*net.OpError); ok {
			if err, ok := err.Err.(net.Error); ok && err.Timeout() {
				return true
			}
		}
	}

	return false
}

func (self *RestClient) GetFullUrl(service string, path string) string {
	var fullUrl string
	if domain, ok := self.config["domain-"+service].(string); ok || "" == service {
		fullUrl = domain + path
	}
	return fullUrl
}

func collectErrors(errs []error) error {
	errMsgs := []string{}
	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}
	err := errors.New(strings.Join(errMsgs, ", "))

	return err
}

func (self *RestClient) setHeaders(ctx context.Context, httpClient *gorequest.SuperAgent, pHeaders *map[string]string) {
	requestId := util.ExtractRequestIDFromCtx(ctx)
	if requestId != util.DefaultRequestID {
		httpClient.Set(util.RequestIDKey, requestId)
		// TODO remove x-req-id
		httpClient.Set("x-req-id", requestId)
	}
	tracingHeader := util.ExtractTracingHeaderFromCtx(ctx)
	if tracingHeader != "" {
		httpClient.Set(util.TracingHeaderKey, tracingHeader)
	}
	if pHeaders != nil {
		for key, value := range *pHeaders {
			httpClient.Set(key, value)
		}
	}
}

func (self *RestClient) RequestXml(ctx context.Context, method string, fullUrl string, data interface{}, pHeaders *map[string]string, params *url.Values) ([]byte, gorequest.Response, error) {
	var paramStr string = ""
	var xmlData []byte
	var err error

	if method != "get" && method != "GET" {
		xmlData, err = xml.Marshal(data)
		if nil != err {
			log.Panic(ctx, "Fail to convert struct to xml", log.Fields{
				"data":     data,
				"error":    err,
				"errorMsg": err.Error(),
			}, err)
		}
	}

	httpClient := self.copyHttpClient(ctx)
	switch method {
	case "GET", "get":
		httpClient.Get(fullUrl)
	case "POST", "post":
		httpClient.Post(fullUrl).Type("xml").Send(string(xmlData))
	case "DELETE", "delete":
		httpClient.Delete(fullUrl).Type("xml").Send(string(xmlData))
	case "PUT", "put":
		httpClient.Put(fullUrl).Type("xml").Send(string(xmlData))
	}

	if pHeaders == nil {
		pHeaders = &map[string]string{}
	}

	(*pHeaders)["Accept"] = "application/xml;charset=utf-8"
	(*pHeaders)["Content-Type"] = "application/xml;charset=utf-8"

	//Form a query-string in url of GET method or body of POST method
	if params != nil {
		paramStr = params.Encode()
		httpClient.Query(paramStr)
	}

	//Add cookies if has
	if len(self.cookies) > 0 {
		httpClient.AddCookies(self.cookies)
	}

	//Support add headers
	self.setHeaders(ctx, httpClient, pHeaders)

	// used to record the duration of request
	// so need to initialize it before calling httpClient.End()
	startTime := time.Now()
	u, err := url.Parse(fullUrl)
	if err != nil {
		u = &url.URL{}
	}
	urlPath := u.Path
	if params != nil {
		encoded := params.Encode()
		if encoded != "" {
			urlPath = urlPath + "?" + encoded
		}
	}
	// sometimes we'll add parameter in fullUrl
	urlParams := u.Query().Encode()
	if urlParams != "" {
		if strings.Contains(urlPath, "?") {
			urlPath = urlPath + "&" + urlParams
		} else {
			urlPath = urlPath + "?" + urlParams
		}
	}
	accessOutLog := log.NewAccessOutLog(
		util.ExtractRequestIDFromCtx(ctx),
		u.Host,
		u.Scheme,
		strings.ToUpper(method),
		urlPath,
		string(xmlData),
		util.GetAccountId(ctx),
	)

	resp, body, errs := httpClient.End()
	statusCode := 599 // default code for request that have no response
	if resp != nil {
		statusCode = resp.StatusCode
	}
	accessOutLog.End(ctx, statusCode, self.RemoteAddr, body, startTime, collectErrors(errs))

	if errs != nil {
		if isRetryable(ctx, errs) {
			errs = retry(RETRY_COUNT, RETRY_TIME_INCR, func() []error {
				startTime := time.Now()
				resp, body, errs = httpClient.End()
				statusCode := 599 // default code for request that have no response
				if resp != nil {
					statusCode = resp.StatusCode
				}
				accessOutLog.End(ctx, statusCode, self.RemoteAddr, body, startTime, collectErrors(errs))
				return errs
			}, isNetTimeout)
		}

		if errs != nil {
			return []byte{}, resp, collectErrors(errs)
		}
	}

	bytes := ([]byte)(body)
	return bytes, resp, err
}

func CtxWithRequestTimeout(ctx context.Context, seconds int64) context.Context {
	return context.WithValue(ctx, REQUEST_TIMEOUT, seconds)
}

func CtxWithBounceToRawString(ctx context.Context, val bool) context.Context {
	return context.WithValue(ctx, BOUNCE_TO_RAW_STRING, val)
}
