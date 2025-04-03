package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/taybart/args"
	"github.com/taybart/env/scan"
)

var (
	app = args.App{
		Name:    "scanenv",
		Version: "v0.0.1",
		Author:  "Taylor Bartlett <taybart@email.com>",
		About:   "check for defined env vars in a project or file",
		Args: map[string]*args.Arg{
			"files": {
				Short: "f",
				Help:  "Comma seperated files to check (./main.go,./util.go)",
			},
			"directory": {
				Short: "d",
				Help:  "Scan directory",
			},
			"validate": {
				Short: "v",
				Help:  "File to validate env config against",
			},
			"print": {
				Short:   "p",
				Help:    "Print contents in env file format, will add file tags above each env",
				Default: false,
			},
			"tags": {
				Short:   "t",
				Help:    "Use go build tags",
				Default: "",
			},
		},
	}
)

func main() {
	if err := app.Parse(); err != nil {
		if errors.Is(err, args.ErrUsageRequested) {
			return
		}
		panic(err)
	}

	var config scan.Config
	if err := app.Marshal(&config); err != nil {
		panic(err)
	}

	res, err := scan.Scan(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(res)
}
