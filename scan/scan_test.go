package scan_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/taybart/env/scan"
)

func TestScan(t *testing.T) {
	is := is.New(t)
	os.Args = []string{"./test", "-d", "./test_project", "-t", "test_tags,other_test"}
	err := scan.Args.Parse()
	is.NoErr(err)
	// hijack stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err = scan.Scan(scan.Args)
	is.NoErr(err)
	w.Close()
	out, _ := io.ReadAll(r)
	// use hijacked output
	os.Stdout = rescueStdout
	is.True(strings.Compare(strings.ReplaceAll(string(out), "\n", ""), `ENV=""PORT="6969"BUILD_TAG=""`) == 0)
}
