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
  os.Setenv("WOOT", `{ "yes": "even_json" }`)

  env.Set([]string{
      "ENVS_ARE_FUN",
      "WOOT",
      "PORT=8080", // default values
      "INSECURE?", // optional values, default to go "zero values"
  })


  config := struct{
    Yes string `json:"yes"`
  }{}

  env.JSON("WOOT", &config)

  fmt.Println(config.Yes)

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

`go install github.com/taybart/env/cmd/scanenv`

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
$ scanenv -d ./ -validate .env
[ERROR] Missing required env
ENV
[WARN] Using default value for PORT=69
```
