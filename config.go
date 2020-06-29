package envconfig

import (
	"errors"
	"fmt"
	"github.com/arturmartini/extjson"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	config = "envconfig"
)

type Configuration struct {
	Envs     []string
	Args     []string
	Required []string
	Default  map[string]string
}

func Initialize(path string, config *Configuration) error {
	if config == nil {
		config = &Configuration{}
	}
	cleanup()
	err := initConfiguration(path, config)
	if err != nil {
		return err
	}

	err = execValidate()
	if err != nil {
		return err
	}

	return nil
}

func GetStr(key string) string {
	return extjson.GetStr(key)
}

func GetInt(key string) int {
	return extjson.GetInt(key)
}

func GetFloat(key string) float64 {
	return extjson.GetFloat(key)
}

func GetListStr(key string) []string {
	return extjson.GetList(key)
}

func GetMapStr(key string) map[string]string {
	return extjson.GetMapStr(key)
}

func execValidate() error {
	errors := []error{}

	err := validateConfig()
	if err != nil {
		errors = append(errors, err)
	}

	return checkError(errors)
}

func load(path, key string) error {
	err := extjson.LoadFile(path, key)
	if err != nil {
		message := fmt.Sprintf("envconfig: error load file path: %s", path)
		log.WithError(err).Warn(message)
		return errors.New(message)
	}
	return nil
}

func initConfiguration(path string, config *Configuration) error {
	errors := []error{}
	err := loadConfig(path)
	if err != nil {
		errors = append(errors, err)
	}

	loadDefault(config)

	err = loadArgs(config)
	if err != nil {
		errors = append(errors, err)
	}

	err = loadEnvs(config)
	if err != nil {
		errors = append(errors, err)
	}

	return checkError(errors)
}

func loadConfig(path string) error {
	err := load(path, config+generateHash())
	if err != nil {
		message := fmt.Sprintf("envconfig: config not detected in path: %s", path)
		log.WithError(err).Warnf(message)
		return errors.New(message)
	}
	return nil
}

func loadDefault(config *Configuration){
	values := extjson.GetMapStr("envconfig.default")
	for k, v := range config.Default {
		if values == nil {
			values = map[string]string{}
		}
		values[k] = v
	}

	ec := extjson.GetMap("envconfig")
	for k, v := range values {
		if ec == nil {
			ec = map[string]interface{}{}
		}
		ec[k] = v
	}

	extjson.Add(ec)
}

func initBaseMap() (map[string]interface{}, map[string]interface{}) {
	envsParams := extjson.GetMap("envconfig")
	envsInnerParams := map[string]interface{}{}
	if len(envsParams) == 0 {
		envsParams["envconfig"] = envsInnerParams
	} else {
		envsInnerParams = envsParams
	}
	return envsParams, envsInnerParams
}

func cleanup() {
	extjson.Cleanup()
}

func generateHash() string {
	return strconv.Itoa(time.Now().Second())
}
