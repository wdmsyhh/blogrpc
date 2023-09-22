package server

import (
	"fmt"
	"log"
	"strings"
	"time"

	"blogrpc/core/extension"
	"blogrpc/openapi/business/controller"
	"blogrpc/openapi/business/middleware"
	"blogrpc/openapi/business/router"
	"blogrpc/openapi/business/util"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/cast"
	conf "github.com/spf13/viper"
)

func Bootstrap(server *ApiServer) error {
	log.Printf("Start to boot server, %+v", server)

	// Load extensions for openapi
	extensions := []string{"request", "mgo", "redis"}
	//extensions := []string{"request", "mgo"}
	extension.LoadExtensionsByName(extensions, server.DebugMode())

	// Init and load system wide middlewares
	initDefaultMiddlewares(server)

	// Load all middlewares
	loadEnabledMiddlewares(server)

	loadAuthMiddleware(server)

	responseWriterMiddleware := &middleware.ResponseWriterMiddleware{}
	server.Engine.Use(responseWriterMiddleware.MiddlewareFunc())
	// Load all controller/action with authentication to server
	controllerActions := controller.GetAllControllerAction()
	for _, a := range controllerActions {
		path := fmt.Sprintf("/%s", strings.Trim(a.Path, "/"))
		a.Method = strings.ToUpper(a.Method)
		if a.FixPath {
			// API path conflict with gin router, use custom router to handle gin NotFound handle
			// Fix: https://github.com/gin-gonic/gin/issues/388
			// Fix: https://github.com/julienschmidt/httprouter/issues/73
			// If httprouter upgrade to v2, maybe this can be solved
			router.Handle(a.Method, path, a)
		} else {
			server.Engine.Handle(a.Method, path, a.GenerateHandler())
		}
	}

	server.Engine.Any("/v2/*rest", server.GenerateV2HandlerFunc())
	// Add the status endpoint handler to show current API status
	server.Engine.GET("/.apistatus", server.GetServerStatus)
	server.Engine.GET("/ping", server.Ping)
	server.Engine.NoRoute(server.NotFound)
	server.Engine.NoMethod(server.MethodNotAllowed)
	return nil
}

func Unload() {
	// Unload extensions
	exts := extension.GetExtensions()
	log.Printf("Start unloading extensions...")

	for name, ext := range exts {
		log.Printf("Unloading ext: %s...", name)
		if nil != ext {
			ext.Close()
			exts[name] = nil
		}
	}
}

func initDefaultMiddlewares(server *ApiServer) {
	server.Engine.Use(middleware.Recovery())
	// Add X_QCRM_TRACE_ID to header
	rewrite := &middleware.RewriteMiddleware{}
	server.Engine.Use(rewrite.MiddlewareFunc())

	statusMiddleware := &middleware.StatusMiddleware{}
	server.Engine.Use(statusMiddleware.MiddlewareFunc())
	go func(sm *middleware.StatusMiddleware) {
		server.ServerStatus = sm.GetStatus()
		// Collect server status in every 10 seconds
		tick := time.Tick(time.Second * 10)
		for {
			select {
			case <-tick:
				server.ServerStatus = sm.GetStatus()
			}
		}
	}(statusMiddleware)

	accessLogJsonMiddleware := &middleware.AccessLogJsonMiddleware{}
	server.Engine.Use(accessLogJsonMiddleware.MiddlewareFunc())
}

func loadAuthMiddleware(server *ApiServer) {
	am := &middleware.AuthMiddleware{
		Version: server.Version,
		Env:     server.Env,
		LookupFunc: func(t *jwt.Token) (interface{}, error) {
			var (
				accountId  string
				secret     string
				err        error
				value      interface{}
				isOldToken bool
				claims     = t.Claims.(jwt.MapClaims)
			)

			// uid only exists in old accessToken's payload(aka Claims)
			if value, isOldToken = claims["uid"]; !isOldToken {
				value = claims["aid"]
			}
			accountId = cast.ToString(value)

			secret, err = getSecret(accountId)
			if err != nil {
				return nil, err
			}

			return []byte(secret), nil
		},
		SkipReferrerValidator: server.DebugMode() || server.DevMode(),
	}

	server.Engine.Use(am.Auth())
}

func loadEnabledMiddlewares(server *ApiServer) {
	ms := middleware.GetMiddlewares()
	configedM := []string{
		"jsonp",
		"cors",
	}
	log.Printf("Start loading and enabling middlewares: %v", configedM)

	// Iterate all configured middlewares and initialize configuration
	for _, m := range configedM {

		if minstance, ok := ms[m]; ok {
			mconfig := conf.GetStringMap(fmt.Sprintf("%s-%s", "middleware", m))
			mconfig["service"] = conf.GetString("service")
			mconfig["addr"] = conf.GetString("addr")
			mconfig["env"] = conf.GetString("env")
			err := minstance.InitWithConf(mconfig, server.DebugMode())

			if err != nil {
				panic(err)
			} else {
				server.Engine.Use(minstance.MiddlewareFunc())
			}
		}
	}
}

func getSecret(accountId string) (string, error) {
	var (
		key    = util.SECRET_HASH_PREFIX + accountId
		secret string
	)

	secret, err := extension.RedisClient.Get(key)
	if err != nil {
		return "", err
	}

	return secret, nil
}
