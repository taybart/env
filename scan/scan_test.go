package scan_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/taybart/env/scan"
)

func TestScan(t *testing.T) {
	is := is.New(t)
	os.Args = []string{"./test", "-f", "./test_project/main.go"}
	err := scan.Args.Parse()
	is.NoErr(err)
	// hijack stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err = scan.Scan(scan.Args)
	is.NoErr(err)
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	// use hijacked output
	is.True(strings.Compare(strings.ReplaceAll(string(out), "\n", ""), `ENV=""PORT="6969"`) == 0)
}
