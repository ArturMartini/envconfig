# Golang Initializer Library (GIL)
GIL is a library to simplify initialize golang program based on json file and variable environments

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
import "github.com/arturmartini/gil"

func main() {
    //Process read a config file and this file can
    //have some metadata to GIL execute validation 
    gil.Initialize("your_path_config_file.json")
    doSomething(gil.GetStr("key"))
}

func doSomething(param string) {
    //you maybe need get webservice, database address or system property
}
```

Validation:
```json
  {
    "foo": "",
    "gil_config": {
      "required": [
        "address"
      ]
    }
  }
```

```go
package main
import (
"fmt"
"github.com/arturmartini/gil"
)

func main() {
    err := gil.Initialize("with_validation.json")
    fmt.Println(err.Error()) //error validate config param required fields: [address]
}
```

Envs:
```json
 {
  "gil_env": {
      "values": {
        "http-port": "8080",
        "address": ""
      },
      "required": [
        "address"
      ]
    }
 }
```

```go
package main
import (
"fmt"
"github.com/arturmartini/gil"
)

func main() {
    //go run main.go address="localhost:8081"
    //When load envs if http-port not found, this map value is set how to default
    //Note this address is required but we pass address by arguments
    //So initialize not return error
    gil.Initialize("with_envs.json")
    fmt.Println(gil.GetStr("http-port")) //"8080"
    fmt.Println(gil.GetStr("address"))   //"localhost:8081"
}
```