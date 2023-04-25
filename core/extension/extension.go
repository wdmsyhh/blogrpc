package extension

import (
	"fmt"
	conf "github.com/spf13/viper"
	"log"
	"os"
)

type IExtension interface {
	Name() string
	InitWithConf(conf map[string]interface{}, debug bool) error
	Close()
}

var (
	//init empty extensions map
	extensions map[string]IExtension = make(map[string]IExtension)
)

func GetExtensions() map[string]IExtension {
	return extensions
}

// GetExtension is used to get the extension instance with name
func GetExtension(name string) IExtension {
	return extensions[name]
}

func LoadEnabledExtensions(enabled map[string]IExtension, isDebug bool) {
	// Iterate enabled extensions and initialize configuration
	for name, ext := range enabled {
		extConfig := conf.GetStringMap(fmt.Sprintf("%s-%s", "extension", name))

		err := ext.InitWithConf(extConfig, isDebug)

		if err != nil {
			log.Printf("[Error] Fail to init extension %s", name)
			panic(err)
		}
	}
}

func LoadExtensions(isDebug bool) {
	exts := GetExtensions()
	log.Println("Start to initialize extensions")

	LoadEnabledExtensions(exts, isDebug)
}

func LoadExtensionsByName(extensions []string, isDebug bool) {
	exts := GetExtensions()
	log.Println("Start to initialize extensions")

	existExts := map[string]IExtension{}
	for _, name := range extensions {
		if e, ok := exts[name]; ok {
			existExts[name] = e
		}
	}

	LoadEnabledExtensions(existExts, isDebug)
}

func RegisterExtension(ext IExtension) {
	if nil != ext {
		name := ext.Name()
		extensions[name] = ext
	} else {
		log.Println("[Error] Failed to load extention, extension is nil.")
	}
}

func getLogger(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix, log.LstdFlags)
}
