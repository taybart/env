# env

[![test](https://github.com/taybart/env/actions/workflows/test.yaml/badge.svg)](https://github.com/taybart/env/actions/workflows/test.yaml)

Easy environments in go!

```go
package main

import (
  "fmt"
  "os"
  "github.com/taybart/env"
)

func main() {
  // set some env
  os.Setenv("ENVS_ARE_FUN", "true")
  os.Setenv("WOOT", `{ "yes": "even_json" }`)
  // Declare our env for this file
  env.Add([]string{
      "ENVS_ARE_FUN",
      "WOOT",
      "PORT=8080", // default values
      "INSECURE?", // optional values, default to go "zero values"
  })


  var config struct{
    Yes string `json:"yes"`
  }

  env.JSON("WOOT", &config)

  fmt.Println(config.Yes)

  // check if vars are defined
  if !env.Has("INSECURE") {
    fmt.Println("INSECURE is not defined")
  }
  // default zero values
  if !env.Bool("INSECURE") {
    fmt.Println("This is super secure now")
  }

  if env.Bool("ENVS_ARE_FUN") {
    fmt.Println("They really are...")
  }

  // look up random envs
  home := env.Get("HOME")
  fmt.Printf("HOME=%s\n", home)
}
```

## Generate env requirements with the CLI

#### Installation

`go install github.com/taybart/env/cmd/scanenv@latest`

#### Generate env file

Single File:

```sh
$ scanenv -f ./main.go
ENVS_ARE_FUN=""
ENVS_ARE_COMPLICATED=""
PORT="8080"
INSECURE="Value marked as optional"
```

Recursive Directory:

```sh
$ scanenv -d .
ENVS_ARE_FUN=""
ENVS_ARE_COMPLICATED=""
PORT="8080"
INSECURE="Value marked as optional"
```

Build tags:

```sh
$ scanenv -d . -t production,fancy
ENVS_ARE_FUN=""
ENVS_ARE_COMPLICATED=""
PORT="8080"
INSECURE="Value marked as optional"
FANCY_PRODUCTION_MARK="✓"
```

#### Get where env is declared

```sh
$ scanenv -d . -p
# ./main.go
ENVS_ARE_FUN=""
ENVS_ARE_COMPLICATED=""
PORT="8080"
INSECURE="Value marked as optional"

# ./util/redis.go
REDIS_PORT=""
```

#### Validate .env file

```sh
$ scanenv -d . --validate .env
[ERROR] Missing required env
ENV
[WARN] Using default value for PORT=69
```
