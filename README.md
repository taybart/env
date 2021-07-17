# env

[![Test](https://github.com/taybart/env/actions/workflows/test.yml/badge.svg)](https://github.com/taybart/env/actions/workflows/test.yml)

An environment package which will take in a list of keys and the secrets manager from AWS and get the default values for any keys that are not already defined as environment variables.

```go
package main

import (
  "github.com/taybart/env"
)

func main() {
  os.Setenv("ENVS_ARE_FUN", "true")
  os.Setenv("ENVS_ARE_COMPLICATED", `{ "test": "json?" }`)

  env.Set([]string{
      "ENVS_ARE_FUN",
      "ENVS_ARE_COMPLICATED",
      "PORT=8080", // default values
      "INSECURE?", // optional values, default to go "zero values"
  })


  config := struct{
    Test string `json:"test"`
  }{}

  env.JSON("ENVS_ARE_COMPLICATED", &config)

  fmt.Println(config.Test)

  // check if vars are defined
  if env.Has("INSECURE") || env.Bool("INSECURE") {
    fmt.Println("This is super insecure now")
  }

  if env.Bool("ENVS_ARE_FUN") {
    fmt.Println("They really are...")
  }

  // look up random envs
  home := env.Get("HOME")
  fmt.Println("HOME=%s", home)
}
```

## CLI

#### Installation

`go install github.com/taybart/environment/cmd/env-scanner`

#### Generate env file

Single File:

```sh
env_scanner -f ./main.go

ENVS_ARE_FUN=""
ENVS_ARE_COMPLICATED=""
PORT="8080"
INSECURE="Value marked as optional"
```

Directory:

```sh
$ env_scanner -d ./
ENVS_ARE_FUN=""
ENVS_ARE_COMPLICATED=""
PORT="8080"
INSECURE="Value marked as optional"
```

#### Get where env is declared

```sh
$ env_scanner -d ./ -p
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
$ env_scanner -d ./ -validate .env
~~Missing Required~~
   ENV
~~Using default~~
   PORT=8080
```