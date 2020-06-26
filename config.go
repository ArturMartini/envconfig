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
	gsilConfig = "gsilConfig"
)

type Configuration struct {
	Envs     []string
	Required []string
	Default  map[string]string
}

func Initialize(path string, config *Configuration) error {
	if config == nil {
		config = &Configuration{}
	}
	cleanup()
	initConfiguration(path, config)
	return execValidate()
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
	return extjson.GetMap(key)
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

func initConfiguration(path string, config *Configuration) {
	loadConfig(path)
	loadEnv(config)
}

func loadConfig(path string) error {
	err := load(path, gsilConfig+generateHash())
	if err != nil {
		message := fmt.Sprintf("envconfig: config not detected in path: %s", path)
		log.WithError(err).Warnf(message)
		return errors.New(message)
	}
	return nil
}

func cleanup() {
	envsConfigured = []string{}
	envsRequired = []string{}
	extjson.Cleanup()
}

func generateHash() string {
	return strconv.Itoa(time.Now().Second())
}
