package gil

import (
	"errors"
	"fmt"
	"github.com/ArturMartini/gel"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	gsilConfig = "gsilConfig"
)

func Initialize(path string) error {
	cleanup()
	initConfiguration(path)
	extractValidations()
	return execValidate()
}

func GetStr(key string) string {
	v, ok := envsParams[key]
	if ok {
		return *v
	}
	return gel.GetStr(key)
}

func GetInt(key string) int {
	return gel.GetInt(key)
}

func GetFloat(key string) float64 {
	return gel.GetFloat(key)
}

func GetListStr(key string) []string {
	return gel.GetList(key)
}

func GetMapStr(key string) map[string]string {
	return gel.GetMap(key)
}

func execValidate() error {
	err := validate()
	if err != nil {
		return err
	}
	return nil
}

func load(path, key string) error {
	err := gel.LoadFile(path, key)
	if err != nil {
		message := fmt.Sprintf("gsil: error load file path: %s", path)
		log.WithError(err).Warn(message)
		return errors.New(message)
	}
	return nil
}

func initConfiguration(path string) {
	loadConfig(path)
	loadEnvs()
}

func loadConfig(path string) error {
	err := load(path, gsilConfig + generateHash())
	if err != nil {
		message := fmt.Sprintf("gsil: config not detected in path: %s", path)
		log.WithError(err).Warnf(message)
		return errors.New(message)
	}
	return nil
}

func cleanup() {
	envsParams     = map[string]*string{}
	envsConfigured = map[string]string{}
	envsRequired   = []string{}
	validations = []ValidateType{}
}

func generateHash() string {
	return strconv.Itoa(time.Now().Second())
}