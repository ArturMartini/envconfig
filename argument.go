package envconfig

import (
	"github.com/arturmartini/extjson"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func loadArgs(config *Configuration) error {
	if !extjson.FoundKey("envconfig") && config == nil {
		log.Warn("envconfig: configuration not found")
		return nil
	}
	envsParams, envsInnerParams := initBaseMap()
	argsConfigured := extjson.GetList("envconfig.args")
	argsConfigured = append(argsConfigured, config.Args...)
	loadArgsConfigured(argsConfigured, envsInnerParams)
	loadArgsRequired(config, envsInnerParams)
	loadArgsByOs(envsParams, argsConfigured)

	extjson.Add(envsParams)
	return nil
}

func loadArgsConfigured(argsConfigured []string, envsInnerParams map[string]interface{}) {
	if len(argsConfigured) > 0 {
		envConfigInterface := []interface{}{}
		for _, v := range argsConfigured {
			envConfigInterface = append(envConfigInterface, v)
		}
		envsInnerParams["args"] = envConfigInterface
	}
}

func loadArgsRequired(config *Configuration, envsInnerParams map[string]interface{}) {
	argsRequired := extjson.GetList("envconfig.required")
	argsRequired = append(argsRequired, config.Required...)

	if len(argsRequired) > 0 {
		envReqInterface := []interface{}{}
		for _, v := range argsRequired {
			envReqInterface = append(envReqInterface, v)
		}
		envsInnerParams["required"] = envReqInterface
	}
}

func loadArgsByOs(envsParams map[string]interface{}, argsConfigured []string) {
	if len(argsConfigured) > 0 {
		for _, arg := range os.Args {
			keyValue := strings.Split(arg, "=")
			if len(keyValue) > 1 {
				key := keyValue[0]
				value := keyValue[1]
				for _, k := range argsConfigured {
					if key == k {
						envsParams[key] = value
					}
				}
			}
		}
	}
}
