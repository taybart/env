package scan_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/taybart/env/scan"
)

func TestScan(t *testing.T) {
	is := is.New(t)
	res, err := scan.Scan(scan.Config{
		Dir:  "./test_project/",
		Tags: "test_tags,other_test",
	})
	is.NoErr(err)
	is.True(strings.Compare(strings.ReplaceAll(res, "\n", ""), `ENV=""PORT="6969"BUILD_TAG=""`) == 0)
}
