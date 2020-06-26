package envconfig

import (
	"fmt"
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	os.Args = append(os.Args, "secret-dir=test", "config-dir=test")
	err := Initialize("test/config.json", nil)
	validateTest(t, nil, err)
}

func TestGets(t *testing.T) {
	os.Args = append(os.Args, "address=test")
	err := Initialize("test/config.json", nil)
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
	expected := "envconfig: error validate required fields: [key1 object1.object1-value]\n"
	err := Initialize("test/config-error.json", nil)
	validateTest(t, expected, err.Error())
}

func TestInitializeEnvRequiredError(t *testing.T) {
	os.Args = []string{}
	expected := "envconfig: error validate required fields: [config-dir secret-dir]\n"
	err := Initialize("test/env-error.json", nil)
	validateTest(t, expected, err.Error())
}

func TestInitializeConfigAndEnvRequiredError(t *testing.T) {
	os.Args = []string{}
	expected := "envconfig: error validate required fields: [key1 object2.object1-value config-dir secret-dir]\n"
	err := Initialize("test/config-env-error.json", nil)
	validateTest(t, expected, err.Error())
}

func TestCleanup(t *testing.T) {
	Initialize("test/config.json", nil)
	v := GetStr("key1")
	validateTest(t, "value1", v)
	cleanup()
	v2 := GetStr("key1")
	validateTest(t, "", v2 )
}

func TestEnvDefault(t *testing.T) {
	Initialize("test/config.json", nil)
	v := GetStr("http-port")
	validateTest(t, "8080", v)
}

func TestAsCodeConfigRequired( t *testing.T) {
	expected := "envconfig: error validate required fields: [key1]\n"
	err := Initialize("test/config-as-code.json", &Configuration{
		Required: []string{"key1"},
	})
	validateTest(t, expected, err.Error())
}

func TestAsCodeConfigRequiredComplex( t *testing.T) {
	expected := "envconfig: error validate required fields: [key1 object1.object5-value]\n"
	err := Initialize("test/config-as-code.json", &Configuration{
		Required: []string{"key1", "object1.object5-value"},
	})
	validateTest(t, expected, err.Error())
}

func TestAsCodeConfigEnv( t *testing.T) {
	os.Args = append(os.Args, "address=localhost:8080")
	err := Initialize("test/config-as-code.json", &Configuration{
		Envs:     []string{"address"},
	})

	addr := GetStr("address")
	validateTest(t, "localhost:8080", addr)
	validateTest(t, nil, err)
}

func TestAsCodeConfigEnvDefault( t *testing.T) {
		os.Args = []string{}
	err := Initialize("test/config-as-code.json", &Configuration{
		Default: map[string]string{
			"http-port": "8081",
		},
	})

	port := GetStr("http-port")
	validateTest(t, "8081", port)
	validateTest(t, nil, err)
}

func TestAsCodeConfigEnvOverrideWithDefault( t *testing.T) {
	os.Args = []string{"http-port=3000"}
	err := Initialize("test/config-as-code.json", &Configuration{
		Default: map[string]string{
			"http-port": "8081",
		},
	})

	port := GetStr("http-port")
	validateTest(t, "3000", port)
	validateTest(t, nil, err)
}

func validateTest(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		reportError(t, expected, actual)
	}
}

func reportError(t *testing.T, expected, actual interface{}) {
	t.Fatal(fmt.Sprintf("Test: %s, Expected: %s, Actual: %s", t.Name(), expected, actual))
}