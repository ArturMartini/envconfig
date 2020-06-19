package gil

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestInitialize(t *testing.T) {
	err := Initialize("config/config.json")
	validateTest(t, nil, err)
}

func TestGets(t *testing.T) {
	os.Args = append(os.Args, "address=test")
	err := Initialize("config/config.json")
	validateTest(t, nil, err)

	vEnvStr := GetStr("address")
	vStr := GetStr("key1")
	vInt := GetInt("key_int")
	vFloat := GetFloat("key_float")
	vList := GetListStr("key_list")
	vMap  := GetMapStr("key_map")

	validateTest(t, "value1", vStr)
	validateTest(t, 1, vInt )
	validateTest(t, 2.01, vFloat)
	validateTest(t, 2, len(vList))
	validateTest(t, 1, len(vMap))
	validateTest(t, "test", vEnvStr)
}

func TestInitializeConfigRequiredError(t *testing.T) {
	expected := "gil: error validate config param required fields: [key1 object1.object1-value]"
	err := Initialize("config/config-error.json")
	validateTest(t, expected, strings.TrimSpace(err.Error()))
}

func TestInitializeEnvRequiredError(t *testing.T) {
	expected := "gil: error validate env param required fields: [config-dir secret-dir]"
	expected2 := "gil: error validate env param required fields: [secret-dir config-dir]"
	err := Initialize("config/env-error.json")
	message := strings.TrimSpace(err.Error())
	match := false
	if message == expected || message == expected2 {
		match = true
	}
	validateTest(t, true, match)
}

func TestInitializeConfigAndEnvRequiredError(t *testing.T) {
	expected := "gil: error validate config param required fields: [key1 object2.object1-value]\n" +
		"gil: error validate env param required fields: [config-dir secret-dir]\n"

	err := Initialize("config/config-env-error.json")

	validateTest(t, len(expected), len(err.Error()))
}

func TestCleanup(t *testing.T) {
	Initialize("config/config.json")
	v := GetStr("key1")
	validateTest(t, "value1", v)
	cleanup()
	v2 := GetStr("key1")
	validateTest(t, "", v2 )
}

func TestEnvDefault(t *testing.T) {
	Initialize("config/config.json")
	v := GetStr("http-port")
	validateTest(t, "8080", v)
}

func validateTest(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		reportError(t, expected, actual)
	}
}

func reportError(t *testing.T, expected, actual interface{}) {
	t.Fatal(fmt.Sprintf("Test: %s, Expected: %s, Actual: %s", t.Name(), expected, actual))
}