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
and validate some. Our library can help you to simplify this process.

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


 