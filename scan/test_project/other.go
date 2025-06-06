//go:build test_tags && other_test

package main

import "github.com/taybart/env"

func init() {
	env.Add([]string{"BUILD_TAG_TEST"})
}
