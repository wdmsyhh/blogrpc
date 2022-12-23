# Introduction

[logrus](https://github.com/sirupsen/logrus) is used as the logging package for both openapi and blogrpc project. Hooks folder contains the customized hook can be used by logrus. For the moment, there is only one SLS log hook which is used to log data to aliyun [SLS service](https://www.aliyun.com/product/sls). Package log contains the utility methods, these methods are created to wrap logrus and pass meta data in context.

# Usage

## Init Logger

Initialize logger based on configuration from file and flag variable

```
coreLog.InitLogger(
    conf.GetString("logger-level"),
    *env,
    "openapi",
    conf.GetStringMapString("extension-sls"),
)
```

## Log Guide

Use golang default log to print bootstrap related information, but use the log package we have encapsulated in core project for other cases.

```
import "git.augmentum.com.cn/blogrpc/core/log"
```

As the example below, use proper log level methods to log only **short and meaningful** message, extra fields can be passed by using `log.Fields` type parameters. All the method should pass extra fields to provide useful context related data fields. We only provide 6 methods to log information:

* `Info`: The info level log will be forward to staging environment and dev environment
* `Warn`: The information which is useful for developer noticing, but not a big problem
* `WarnTrace`: Similiar with `Warn` method, but trace as the last extra parameter
* `Error`: Indicate that this is a bug needed to be fixed if occurs
* `ErrorTrace`: Similiar with `Warn` method, but trace as the last extra parameter
* `Panic`: Print error level information and panic with the last parameter

Commonly, you only need to use `Info`, `Warn`, `Error` method to log extra context related information as the example below.

```
log.Error(c, "Fail to push message", log.Fields{
    "url":      url,
    "headers":  headers,
    "body":     messageBody,
    "response": response,
})
```

But you can use `ErrorTrace` if trace information is needed for some cases.

```
log.ErrorTrace(c, err.Msg, log.Fields{
    "code":  err.Code,
    "error": err.InternalError,
}, debug.Stack())
```

**Notice:** context should always be the first parameter to be passed to the log methods in core project, and extra fields is needed for all methods to provide useful extra information

