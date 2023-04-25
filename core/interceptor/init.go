package interceptor

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cast"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Interceptor interface {
	Name() string // The name of specific interceptor
	Handle(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	InitWithConf(conf map[string]interface{}, debug bool) error
}

var (
	allInterceptors map[string]Interceptor = make(map[string]Interceptor)
	// 开启的拦截器，注意顺序
	enabledInterceptorNames []string = []string{
		"Recovery",
		"Logging",
		"FormatError",
	}
)

func AddInterceptor(i Interceptor) {
	if i != nil {
		allInterceptors[i.Name()] = i
	} else {
		log.Panicln("Failed to add nil interceptor")
	}
}

func GetInterceptors() map[string]Interceptor {
	return allInterceptors
}

func GetSortedInterceptors() []Interceptor {
	var sortedInterceptors []Interceptor
	for _, name := range enabledInterceptorNames {
		incp, ok := allInterceptors[name]
		if ok {
			sortedInterceptors = append(sortedInterceptors, incp)
		}
	}
	return sortedInterceptors
}

func LoadInterceptors(conf map[string]interface{}, debug bool) {
	log.Println("Start initializing interceptors")

	// Iterate all available extensions and initialize configuration
	for _, ext := range allInterceptors {
		var err error
		key := fmt.Sprintf("%s-%s", "interceptor", strings.ToLower(ext.Name()))
		config, ok := conf[key]
		if !ok {
			log.Println("Init with no configuration of " + key)
		}

		cf := cast.ToStringMap(config)
		cf["service"] = conf["service"]
		cf["addr"] = conf["addr"]
		cf["env"] = conf["env"]
		err = ext.InitWithConf(cf, debug)

		if err != nil {
			log.Printf("interceptor: %s, configures: %v, Failed to init extension, error with: %v", ext.Name(), config, err)
			panic(err)
		} else {
			log.Printf("interceptor: %s, configures: %v, Loaded extension and initialized..", ext.Name(), config)
		}
	}
}
