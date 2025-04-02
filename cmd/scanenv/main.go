package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/taybart/args"
	"github.com/taybart/env/scan"
)

func main() {

	app := args.App{}
	app = app.Import(scan.Args)
	if err := app.Parse(); err != nil {
		if errors.Is(err, args.ErrUsageRequested) {
			return
		}
		panic(err)
	}

	if err := scan.Scan(app); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
