package gil

import (
	"errors"
	"fmt"
	"github.com/ArturMartini/gel"
	log "github.com/sirupsen/logrus"
)

var (
	validations = []ValidateType{}
	mapValidateToField = map[ValidateType]string{
		ConfigType: "gsil_config.required",
		EnvType:    "gsil_env.required",
	}
	mapErrValidType = map[ValidateType]string{
		ConfigType: "config param",
		EnvType: "env param",
	}
)

const (
	ConfigType ValidateType = 1
	EnvType ValidateType = 2
)

type ValidateType int

func validate() error {
	errors := []error{}
	for _, v := range validations {
		err := executeValidate(v)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return checkError(errors)
}

func executeValidate(validateType ValidateType) error {
	switch validateType {
	case ConfigType:
		return validateConfig()
	case EnvType:
		return validateEnv()
	default:
		return nil
	}
}

func validateEnv() error {
	missing := []string{}
	for k, v := range envsConfigured {
		isRequired := false
		for _, vEnv := range envsRequired {
			if k == vEnv {
				isRequired = true
			}
		}

		if isRequired && v == "" {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		message := fmt.Sprintf("gsil: error validate %s required fields: %v", mapErrValidType[EnvType], missing)
		log.Warn(message)
		return errors.New(message)
	}
	return nil
}

func validateConfig() error {
	missing := []string{}
	requireds := gel.GetList(mapValidateToField[ConfigType])
	for _, k := range requireds {
		value := gel.GetStr(k)
		if value == "" {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		message := fmt.Sprintf("gsil: error validate %s required fields: %v", mapErrValidType[ConfigType], missing)
		log.Warn(message)
		return errors.New(message)
	}
	return nil
}

func extractValidations() {
	validations = []ValidateType{}
	conf := gel.GetList(mapValidateToField[ConfigType])
	if len(conf) > 0 {
		validations = append(validations, ConfigType)
	}

	if len(envsRequired) > 0 {
		validations = append(validations, EnvType)
	}
}

func checkError(errs []error) error {
	message := ""
	for _, e := range errs {
		message += e.Error() + "\n"
	}
	if len(errs) > 0 {
		return errors.New(message)
	}
	return nil
}