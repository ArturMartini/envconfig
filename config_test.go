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

func TestInitializeFileNotFound(t *testing.T) {
	expected := "envconfig: config not detected in path: not.exists\n"
	err := Initialize("not.exists", nil)
	validateTest(t, expected, err.Error())
}

func TestGets(t *testing.T) {
	os.Args = append(os.Args, "address=arg_test")
	os.Setenv("env1", "env1")
	err := Initialize("test/config.json", nil)
	validateTest(t, nil, err)

	vArgStr := GetStr("address")
	vEnvStr := GetStr("env1")
	vStr := GetStr("key1")
	vInt := GetInt("key_int")
	vFloat := GetFloat("key_float")
	vList := GetListStr("key_list")
	vMap := GetMapStr("key_map")
	vMapInterface := GetMap("key_map")
	vBool := GetBool("key_bool")

	validateTest(t, "value1", vStr)
	validateTest(t, 1, vInt)
	validateTest(t, 2.01, vFloat)
	validateTest(t, 2, len(vList))
	validateTest(t, 1, len(vMap))
	validateTest(t, 1, len(vMapInterface))
	validateTest(t, "arg_test", vArgStr)
	validateTest(t, "env1", vEnvStr)
	validateTest(t, true, vBool)
}

func TestGetsAsCode(t *testing.T) {
	os.Args = append(os.Args, "address=arg_test")
	os.Setenv("env_test", "env_test")
	config := &Configuration{
		Envs: []string{
			"env_test",
		},
		Args: []string{
			"address",
		},
		Default: map[string]string{
			"def": "v1",
		},
	}

	err := Initialize("test/config-as-code.json", config)
	validateTest(t, nil, err)

	vArgStr := GetStr("address")
	vEnvStr := GetStr("env_test")
	vStr := GetStr("key")
	vInt := GetInt("key_int")
	vFloat := GetFloat("key_float")
	vList := GetListStr("key_list")
	vMap := GetMapStr("key_map")
	defStr := GetStr("def")

	validateTest(t, "v1", defStr)
	validateTest(t, "value1", vStr)
	validateTest(t, 1, vInt)
	validateTest(t, 2.01, vFloat)
	validateTest(t, 2, len(vList))
	validateTest(t, 1, len(vMap))
	validateTest(t, "arg_test", vArgStr)
	validateTest(t, "env_test", vEnvStr)
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
	validateTest(t, "", v2)
}

func TestEnvDefault(t *testing.T) {
	Initialize("test/config.json", nil)
	v := GetStr("http-port")
	validateTest(t, "8080", v)
}

func TestAsCodeConfigRequired(t *testing.T) {
	expected := "envconfig: error validate required fields: [key1]\n"
	err := Initialize("test/config-as-code.json", &Configuration{
		Required: []string{"key1"},
	})
	validateTest(t, expected, err.Error())
}

func TestAsCodeConfigRequiredComplex(t *testing.T) {
	expected := "envconfig: error validate required fields: [key1 object1.object5-value]\n"
	err := Initialize("test/config-as-code.json", &Configuration{
		Required: []string{"key1", "object1.object5-value"},
	})
	validateTest(t, expected, err.Error())
}

func TestAsCodeConfigArgs(t *testing.T) {
	os.Args = append(os.Args, "address=localhost:8080")
	err := Initialize("test/config-as-code.json", &Configuration{
		Args: []string{"address"},
	})

	addr := GetStr("address")
	validateTest(t, "localhost:8080", addr)
	validateTest(t, nil, err)
}

func TestAsCodeConfigDefault(t *testing.T) {
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

func TestAsCodeConfigDefaultArgsOverrideWithArg(t *testing.T) {
	os.Args = []string{"http-port=3000"}
	err := Initialize("test/config-as-code.json", &Configuration{
		Args: []string{"http-port"},
		Default: map[string]string{
			"http-port": "8081",
		},
	})

	port := GetStr("http-port")
	validateTest(t, "3000", port)
	validateTest(t, nil, err)
}

func TestInitializeWithVariableEnvironment(t *testing.T) {
	os.Setenv("env1", "localhost:8080")
	err := Initialize("test/config.json", nil)
	expected := "localhost:8080"
	addr := GetStr("env1")
	validateTest(t, nil, err)
	validateTest(t, expected, addr)
}

func TestEnvOrverideDefault(t *testing.T) {
	os.Setenv("env_override", "override")
	err := Initialize("test/config.json", &Configuration{
		Envs: []string{"env_override"},
		Default: map[string]string{
			"env_override": "any",
		},
	})
	expected := "override"
	result := GetStr("env_override")
	validateTest(t, nil, err)
	validateTest(t, expected, result)
}

func validateTest(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		reportError(t, expected, actual)
	}
}

func reportError(t *testing.T, expected, actual interface{}) {
	t.Fatal(fmt.Sprintf("Test: %s, Expected: %s, Actual: %s", t.Name(), expected, actual))
}
