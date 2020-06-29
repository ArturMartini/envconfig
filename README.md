  # ENVCONFIG
Envconfig is a library to simplify initialize golang program based on json file, program arguments and environments variable with configuration and validation in embedded file or as code

## Features 
* Initialize properties by json file  
* Initialize properties by variable environment
* Initialize properties by arguments
* Check required properties by json config
* Check required properties by as code
* Set default value
* Abstract the complex bind and validate structure

Commonly systems need get program arguments, environments variable, properties in config file
and validate. Our library can hep you to simplify this process.

Usage:
```go
package main
import "github.com/arturmartini/envconfig"

func main() {
    //Process read a config file and this file can
    //have some metadata to envconfig to check required fields and set default values
    //can send metadata as code with Configuration struct in Initialize func
    //can accesss sub values with separator "."
    envconfig.Initialize("your_path_config_file.json", nil)
    doSomething(envconfig.GetStr("key"))
    doSomething(envconfig.GetStr("obj.key"))
}

func doSomething(param string) {
    //you maybe need get webservice, database address or system property
}
```

Configuration in file:
```json
  {
    "foo": "",
    "envconfig": {
      "required": [
        "address"
      ],
      "envs" :[
        "address",
        "port"
      ],
      "args]" :[
        "arg1" 
      ],
      "default":{
        "port": "8080"
      }
    }
  }
```
* The fields required is properties, variable environment or program arguments and all is validate
* Default to be definied is default value to key, if key has value in file, enviroment or aguments the default value is overrided


```go
package main
import (
  "fmt"
  "github.com/arturmartini/envconfig"
)

func main() {
    err := envconfig.Initialize("with_validation.json", nil)
    fmt.Println(err.Error()) //envconfig: error validate required fields: [address]
}
```

* Validation as code
```go
package main
import (
  "fmt"
  "github.com/arturmartini/envconfig"
)

func main() {
    //export address="localhost"
    //go run main.go port=":8081"
    //When load envs if http-port not found, this map defaut is set to default value
    //Note this address is required but we pass address by arguments
    config := &Configuration{
      Envs: []string{"address"},
      Agrs: []string{"port"}
      Required: []string{"address"},
      Default: map[string]string{
          "defualt","localhost:8080"
      }
    }
    envconfig.Initialize("without_file_validation.json", config)
    fmt.Println(envconfig.GetStr("address"))   //"localhost"
    fmt.Println(envconfig.GetStr("port"))      //":8081"
    fmt.Println(envconfig.GetStr("protocol"))  //"localhost:8080"
}
```