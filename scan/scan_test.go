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
	is.True(res.Equal(scan.Env{
		Values: map[string]scan.EnvVar{
			// main.go
			"ENV":    {},
			"PORT":   {Value: "6969", HasDefault: true},
			"SECURE": {Optional: true},
			// other.go (with build tags)
			"BUILD_TAG_TEST": {},
		}},
	))
	resF := strings.ReplaceAll(res.ToFile(), "\n", "")
	is.True(strings.Compare(resF, `BUILD_TAG_TEST=""ENV=""PORT="6969"SECURE="value is marked as optional"`) == 0)

	// fmt.Println(res.EnvByFile())
}
