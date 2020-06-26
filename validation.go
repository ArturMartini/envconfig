package envconfig

import (
	"errors"
	"fmt"
	"github.com/arturmartini/extjson"
	log "github.com/sirupsen/logrus"
)

func validateConfig() error {
	missing := []string{}
	requireds := extjson.GetList("envconfig.required")
	for _, k := range requireds {
		if !extjson.FoundKey(k) {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		message := fmt.Sprintf("envconfig: error validate required fields: %v", missing)
		log.Warn(message)
		return errors.New(message)
	}
	return nil
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
