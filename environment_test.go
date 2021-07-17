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
	env.Set([]string{fmt.Sprintf("%s=%s", k, v)})
	is.True(env.Is(k, v))
}

// Test that optionals are set to zero value
func TestOptionalKey(t *testing.T) {
	is := is.New(t)

	k := "TEST_OPTIONAL_KEY"

	// Add optional to env
	env.Set([]string{
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

	type I struct {
		Keys map[string]map[string]string `json:"keys"`
	}
	// expected by call to be set, particular value doesn't matter
	os.Setenv(k, `{"keys": {"embedded": "someInnerVal"}, {"someOtherKey": "key": "val"}}`)

	// test struct
	var returned I
	err := env.JSON(k, &returned)
	is.NoErr(err)

	expected := "someInnerVal"
	is.True(returned.Keys["keys"]["embedded"] != expected)
}
