package main

import (
	"fmt"

	"github.com/taybart/env"
)

func main() {
	env.Set([]string{
		"ENV",
		"PORT=69",
	})

	if env.Is("ENV", "production") {
		fmt.Println("HOLY CRAP CALL THE SENIOR")
	}
}
