package env

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/taybart/log"
)

var (
	keyRE = regexp.MustCompile(`([[:word:]]+)([=?])?(.*)?`)
	// Optional keys that should be set to zero value
	optionalKeys map[string]bool
)

/* Add : environment variables for use later. This is global to the project
 * requred -> NAME
 * with_default -> NAME=taybart
 * optional -> NAME? // defaults to zero value
 */
func Add(keys []string) {
	err := Ensure(keys)
	if err != nil {
		panic(err)
	}
}

// Ensure : check that env vars are defined, set default, mark optional
func Ensure(keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	missingKeys := []string{}
	optionalKeys = GetOptional(keys)
	for _, key := range keys {
		fkey, val := GetDefault(key) // formatted key and default value
		_, found := os.LookupEnv(fkey)
		if !found {
			if optionalKeys[fkey] {
				log.Warnf("%s marked optional and not defined\n", fkey)
				continue
			}
			if val != "" { // is there a default value?
				log.Warnf("Setting %s to default value of %s\n", fkey, val)
				os.Setenv(fkey, val)
				continue
			}
			missingKeys = append(missingKeys, key)
		}
	}
	for _, key := range missingKeys {
		log.Errorf("Missing environment variable: %s%s%s\n", log.Red, key, log.Reset)
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("set all required environment variables: %v", missingKeys)
	}
	return nil
}

// Has : see if defined or has optional
func Has(key string) bool {
	_, b := os.LookupEnv(key)
	return b
}

// Is : returns if the variable _is_ the string
func Is(key, compare string) bool {
	if val, found := os.LookupEnv(key); found {
		return val == compare
	}
	return false
}

// Get : returns the environment value as a string
func Get(key string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}
	log.Warnf("getting optional key %v\n", optionalKeys)
	if _, found := optionalKeys[key]; found {
		return ""
	}

	log.Fatal("Trying to retrieve uninitialized environment variable:", key)
	return ""
}

// Decode : returns the environment value as a string decoded base64
func Decode(key string) (string, error) {
	if val, found := os.LookupEnv(key); found {
		decoded, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}
	log.Warnf("getting optional key %v\n", optionalKeys)
	if _, found := optionalKeys[key]; found {
		return "", nil
	}

	log.Fatal("Trying to retrieve/decode uninitialized environment variable:", key)
	return "", nil
}

// Int : returns the key as an int or panics
func Int(key string) int {
	if val, found := os.LookupEnv(key); found {
		converted, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("An error occurred in converting the value [%s] retrieved with key [%s] to an int: %s", val, key, err)
		}
		return converted
	}
	if _, found := optionalKeys[key]; found {
		return 0
	}

	log.Fatal("Trying to retrieve uninitialized environment variable:", key)
	return 0
}

// Bool : returns the env var as its value, or false if it doesn't exist
func Bool(key string) bool {
	if val, found := os.LookupEnv(key); found {
		return val == "true"
	}
	if _, found := optionalKeys[key]; found {
		return false
	}

	log.Fatal("Trying to retrieve uninitialized environment variable:", key)
	return false
}

// IsSet : returns if the environment variable is set
func IsSet(key string) bool {
	_, found := os.LookupEnv(key)
	return found
}

// JSON : returns the environment value marshalled to input
func JSON(key string, input interface{}) error {
	if val, found := os.LookupEnv(key); found {
		err := json.Unmarshal([]byte(val), input)
		if err != nil {
			return fmt.Errorf("could not unmarshal %s: (value: %+v) %v", key, val, err)
		}
		return nil
	}
	if _, found := optionalKeys[key]; found {
		return nil
	}

	log.Fatalf("Trying to retrieve uninitialized environment variable: %s, found (should show nothing) %s\n", key, os.Getenv(key))
	return nil
}

func GetDefault(entry string) (key string, defaultValue string) {
	res := keyRE.FindAllStringSubmatch(entry, -1)
	return res[0][1], res[0][3]
}

func GetOptional(keys []string) map[string]bool {
	optionals := make(map[string]bool)
	if len(keys) == 0 {
		return optionals
	}
	for _, key := range keys {
		res := keyRE.FindAllStringSubmatch(key, -1)
		optionals[res[0][1]] = res[0][2] == "?"
	}
	return optionals
}

// NoWarn : remove warning logs
func NoWarn() {
	log.SetLevel(log.ERROR)
}

// NoLog : disable logging
func NoLog() {
	log.SetLevel(log.NONE)
}
