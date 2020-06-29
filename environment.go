package envconfig

import (
	"github.com/arturmartini/extjson"
	"os"
	"strings"
)

func loadEnvs(config *Configuration) error {
	envsParams, envsInnerParams := initBaseMap()
	envsConfigured := extjson.GetList("envconfig.envs")
	envsConfigured = append(envsConfigured, config.Envs...)

	loadEnvsRequired(config, envsConfigured, envsInnerParams)
	loadEnvironments(envsConfigured, envsParams)
	extjson.Add(envsParams)
	return nil
}

func loadEnvsRequired(config *Configuration, envsConfigured []string, envsInnerParams map[string]interface{}) {
	envsRequired := extjson.GetList("envconfig.required")
	envsRequired = append(envsRequired, config.Required...)
	if len(envsConfigured) > 0 {
		envConfigInterface := []interface{}{}
		for _, v := range envsRequired {
			envConfigInterface = append(envConfigInterface, v)
		}
		envsInnerParams["args"] = envConfigInterface
	}
}

func loadEnvironments(envsConfigured []string, envsParams map[string]interface{}) {
	if len(envsConfigured) > 0 {
		for _, arg := range os.Environ() {
			keyValue := strings.Split(arg, "=")
			if len(keyValue) > 1 {
				key := keyValue[0]
				for _, k := range envsConfigured {
					if key == k {
						value := keyValue[1]
						envsParams[key] = value
					}
				}
			}
		}
	}
}
