  # ENVCONFIG
Envconfig is a library to simplify initialize golang program based on json file and variable environments with configuration and validation in embedded file or as code

## Features 
* Initialize config by json file  
* Initialize variable environment
* Validate config properties by json config
* Validate environments properties by json config
* Set default value to environments properties
* Abstract the complex bind and validate structure

Commonly systems need get variable environments, properties in config file
and validate some. Our library can hep you to simplify this process.

Usage:
```go
package main
import "github.com/arturmartini/envconfig"

func main() {
    //Process read a config file and this file can
    //have some metadata to envconfig execute validation
    //or you can pass metadata as code with Configuration struct in initialize func
    envconfig.Initialize("your_path_config_file.json", nil)
    doSomething(envconfig.GetStr("key"))
}

func doSomething(param string) {
    //you maybe need get webservice, database address or system property
}
```

Congiguration inner file:
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
      "default":{
        "port": "8080"
      }
    }
  }
```
* The fields required is properties and variable environment and both is validate
* Variable environments to be binding from arguments
* Default to be definied is default value to key, if key has value in file or enviroment the default value is overrided


```go
package main
import (
"fmt"
"github.com/arturmartini/envconfig"
)

func main() {
    err := envconfig.Initialize("with_validation.json", nil)
    fmt.Println(err.Error()) //error validate param required fields: [address]
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
    //go run main.go address="localhost:8081"
    //When load envs if http-port not found, this map defaut is set to default value
    //Note this address is required but we pass address by arguments
    //So initialize not return error
    config := &Configuration{
      Envs: []string{"http-port", "address"},
      Required: []string{"address"},
      Default: map[string]string{
          "http-port","8080"
      }
    }
    envconfig.Initialize("without_file_validation.json", config)
    fmt.Println(envconfig.GetStr("http-port")) //"8080"
    fmt.Println(envconfig.GetStr("address"))   //"localhost:8081"
}
```