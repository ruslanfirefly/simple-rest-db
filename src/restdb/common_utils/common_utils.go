package common_utils

import (
	"github.com/ivahaev/go-logger"
)

func ErrorHandler(v interface{}) {
	if v != nil {
		logger.Error(v)
	}
}

func ParseMap(m map[interface{}]interface{}) map[string]interface{} {
	var resultMap = make(map[string]interface{})

	for key, value := range m {
		switch key := key.(type) {
		case string:
			switch value := value.(type) {
			case string:
				resultMap[key] = value
			case map[interface{}]interface{}:
				resultMap[key] = ParseMap(value)
			}
		}
		//if m[k].(Type) == map[interface{}]interface{})  {
		//	resultMap[reflect.ValueOf(k).String()] = ParseMap()
		//}else {
		//	resultMap[reflect.ValueOf(k).String()] = m[k]
		//}
	}
	return resultMap
}
