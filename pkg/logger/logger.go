package logger

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
)

const moduleKey = "module"

func GetLogger(x any) *Adapter {
	if v, ok := x.(string); ok {
		return &Adapter{
			Logger: log.Logger.With().Str(moduleKey, strings.ToLower(v)).Logger(),
		}
	}

	val := reflect.ValueOf(x)
	if val.Kind() == reflect.Func {
		runtimeFunc := runtime.FuncForPC(val.Pointer())
		if runtimeFunc != nil {
			moduleNames := strings.Split(runtimeFunc.Name(), ".")
			if len(moduleNames) > 1 {
				lastPart := moduleNames[len(moduleNames)-1]
				moduleName := strings.Split(lastPart, "/")
				if len(moduleName) > 0 {
					return &Adapter{
						Logger: log.Logger.With().Str(moduleKey, strings.ToLower(moduleName[len(moduleName)-1])).Logger(),
					}
				}
			}
		}
	}

	return &Adapter{
		Logger: log.Logger.With().Str(moduleKey, "unknown").Logger(),
	}
}
