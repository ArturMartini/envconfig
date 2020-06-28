package envconfig

import (
	"github.com/arturmartini/extjson"
	"os"
	"strings"
)

func loadEnvs(config *Configuration) error {
	envsConfigured := extjson.GetList("envconfig.envs")
	envsConfigured = append(envsConfigured, config.Envs...)
	envsParams := extjson.GetMap("envconfig")
	envsInnerParams := map[string]interface{}{}
	if envsParams == nil {
		envsParams = map[string]interface{}{
			"envconfig": envsInnerParams,
		}
	} else {
		envsInnerParams = envsParams
	}

	fieldsRequired := extjson.GetList("envconfig.required")
	if len(envsConfigured) > 0 {
		envConfigInterface := []interface{}{}
		for _, v := range fieldsRequired {
			envConfigInterface = append(envConfigInterface, v)
		}
		envsInnerParams["args"] = envConfigInterface

	}

	for _, arg := range os.Environ() {
		keyValue := strings.Split(arg, "=")
		if len(keyValue) > 1 {
			key := keyValue[0]
			for _, k := range envsConfigured {
				if key != k {
					continue
				}
			}
			value := keyValue[1]
			envsParams[key] = value
		}
	}
	extjson.Add(envsParams)

	return nil
}
