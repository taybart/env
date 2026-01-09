package env_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/taybart/env"
)

func TestDefault(t *testing.T) {
	is := is.New(t)
	// Define key
	k := "TEST_DEFAULT"
	v := "default_value"
	env.Add([]string{fmt.Sprintf("%s=%s", k, v)})
	is.True(env.Is(k, v))
}

func TestDefaultGuard(t *testing.T) {
	is := is.New(t)
	defer func() {
		// we should panic here
		is.True(recover() != nil)
	}()
	// Define key
	k := "TEST_DEFAULT_GUARD"
	env.Add([]string{fmt.Sprintf("%s=1", k)})
	env.Add([]string{fmt.Sprintf("%s=2", k)})
}

// Test that optionals are set to zero value
func TestOptionalKey(t *testing.T) {
	is := is.New(t)

	k := "TEST_OPTIONAL_KEY"

	// Add optional to env
	env.Add([]string{
		fmt.Sprintf("%s?", k),
	})
	// make sure its empty string
	is.True(env.Get(k) == "")
	// make sure its false
	is.True(!env.Bool(k))
	// make sure its zero
	is.True(env.Int(k) == 0)
}

func TestGet(t *testing.T) {
	is := is.New(t)
	k := "TestGet"
	// set var
	os.Setenv(k, "cool variable")
	// Should return true since TESTING_ENV is set to true
	is.True(env.Get(k) == "cool variable")
}

func TestDecode(t *testing.T) {
	is := is.New(t)
	k := "TestGet"
	// set var
	os.Setenv(k, "Y29vbCB2YXJpYWJsZQ==")
	val, err := env.Decode(k)
	is.NoErr(err)
	is.True(string(val) == "cool variable")
}

// TestHas : if value is set env.Has returns true
func TestHas(t *testing.T) {
	is := is.New(t)
	// Define key
	key := "TEST_HAS"

	// set env
	os.Setenv(key, "this is defined now")
	// set
	is.True(env.Has(key))
}

func TestBool(t *testing.T) {
	is := is.New(t)

	// Define key
	k := "TEST_BOOL"

	// Set env
	os.Setenv(k, "true")
	// Should return true since TESTING_ENV is set to true
	is.True(env.Bool(k))
}

func TestIs(t *testing.T) {
	is := is.New(t)
	os.Setenv("TEST_IS", "testing")
	// Set
	is.True(env.Is("TEST_IS", "testing"))
}

// Test json interface marshaling
func TestInterface(t *testing.T) {
	is := is.New(t)

	// Define key
	k := "TEST_INTERFACE"

	// expected by call to be set, particular value doesn't matter as long as the type is correct
	os.Setenv(k, `{"key": "val", "other": "sudo su"}`)

	// test struct
	var returned map[string]string
	is.NoErr(env.JSON(k, &returned))

	// Should get the correct value
	is.True(returned["key"] == "val")
}
