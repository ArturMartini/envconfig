package envconfig

import (
	"github.com/arturmartini/extjson"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	argsConfigured = []string{}
	argsRequired   = []string{}
)

func loadArgs(config *Configuration) error {
	if !extjson.FoundKey("envconfig") && config == nil {
		log.Warn("envconfig: configuration not found")
		return nil
	}

	argsConfigured = extjson.GetList("envconfig.args")
	argsConfigured = append(argsConfigured, config.Args...)
	envsParams := extjson.GetMap("envconfig")
	envsInnerParams := map[string]interface{}{}
	if envsParams == nil {
		envsParams = map[string]interface{}{
			"envconfig": envsInnerParams,
		}
	} else {
		envsInnerParams = envsParams
	}

	if len(argsConfigured) > 0 {
		envConfigInterface := []interface{}{}
		for _, v := range argsConfigured {
			envConfigInterface = append(envConfigInterface, v)
		}
		envsInnerParams["args"] = envConfigInterface

	}

	argsRequired = extjson.GetList("envconfig.required")
	argsRequired = append(argsRequired, config.Required...)

	if len(argsRequired) > 0 {
		envReqInterface := []interface{}{}
		for _, v := range argsRequired {
			envReqInterface = append(envReqInterface, v)
		}
		envsInnerParams["required"] = envReqInterface
	}

	for _, arg := range os.Args {
		keyValue := strings.Split(arg, "=")
		if len(keyValue) > 1 {
			key := keyValue[0]
			value := keyValue[1]
			envsParams[key] = value

		}
	}

	extjson.Add(envsParams)
	return nil
}
