package main

import (
	"fmt"
	"os"

	"github.com/taybart/args"
	"github.com/taybart/env/scan"
)

func main() {
	app := args.App{}
	app = app.Import(scan.Args)
	app.Parse()

	if err := scan.Scan(app); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
