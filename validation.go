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
		configType: "gil_config.required",
		envType:    "gil_env.required",
	}
	mapErrValidType = map[ValidateType]string{
		configType: "config param",
		envType: "env param",
	}
)

const (
	configType ValidateType = 1
	envType ValidateType = 2
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
	case configType:
		return validateConfig()
	case envType:
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
		message := fmt.Sprintf("gil: error validate %s required fields: %v", mapErrValidType[envType], missing)
		log.Warn(message)
		return errors.New(message)
	}
	return nil
}

func validateConfig() error {
	missing := []string{}
	requireds := gel.GetList(mapValidateToField[configType])
	for _, k := range requireds {
		value := gel.GetStr(k)
		if value == "" {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		message := fmt.Sprintf("gil: error validate %s required fields: %v", mapErrValidType[configType], missing)
		log.Warn(message)
		return errors.New(message)
	}
	return nil
}

func extractValidations() {
	validations = []ValidateType{}
	conf := gel.GetList(mapValidateToField[configType])
	if len(conf) > 0 {
		validations = append(validations, configType)
	}

	if len(envsRequired) > 0 {
		validations = append(validations, envType)
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