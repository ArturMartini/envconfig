package envconfig

import (
	"github.com/arturmartini/extjson"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	envsDefault    = map[string]string{}
	envsConfigured = []string{}
	envsRequired   = []string{}
)

func loadEnv(config *Configuration) error {
	envsConfigured = extjson.GetList("envconfig.envs")
	envsConfigured = append(envsConfigured, config.Envs...)
	mapEnvConfig   := map[string]interface{}{}
	envsParams := map[string]interface{}{
		"envconfig": mapEnvConfig,
	}

	if len(envsConfigured) > 0 {
		envConfigInterface := []interface{}{}
		for _, v := range envsRequired {
			envConfigInterface = append(envConfigInterface, v)
		}
		mapEnvConfig["envs"] = envConfigInterface

	}

	envsRequired = extjson.GetList("envconfig.required")
	envsRequired = append(envsRequired, config.Required...)

	if len(envsRequired) > 0 {
		envReqInterface := []interface{}{}
		for _, v := range envsRequired {
			envReqInterface = append(envReqInterface, v)
		}
		mapEnvConfig["required"] = envReqInterface
	}

	envsDefault = extjson.GetMap("envconfig.default")
	for k, v := range config.Default {
		if len(envsDefault) <= 0 {
			envsDefault = map[string]string{}
		}
		envsDefault[k] = v
	}

	if len(envsDefault) > 0 {
		mapEnvConfig["default"] = envsDefault
	}

	if !extjson.FoundKey("envconfig") && config == nil {
		log.Warn("envconfig: configuration not found")
	}

	for k, v := range envsDefault {
		envsParams[k] = v
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
