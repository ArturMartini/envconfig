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
    envconfig.Initialize("your_path_config_file.json")
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
    err := envconfig.Initialize("with_validation.json")
    fmt.Println(err.Error()) //error validate param required fields: [address]
}
```

Envs:
```json
 {
  "envconfig": {
      "envs": [
        "http-port",
        "address"
      ],
      "required": [
        "address"
      ],
      "default":{
        "http-port": "8080"
      }
    }
 }
```

```go
package main
import (
"fmt"
"github.com/arturmartini/envconfig"
)

func main() {
    //go run main.go address="localhost:8081"
    //When load envs if http-port not found, this map value is set how to default
    //Note this address is required but we pass address by arguments
    //So initialize not return error
    envconfig.Initialize("with_envs.json")
    fmt.Println(envconfig.GetStr("http-port")) //"8080"
    fmt.Println(envconfig.GetStr("address"))   //"localhost:8081"
}
```