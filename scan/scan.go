package scan

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/taybart/args"
	"github.com/taybart/env"
	"github.com/taybart/log"
)

var (
	Args = args.App{
		Name:    "scanenv",
		Version: "v0.0.1",
		Author:  "Taylor Bartlett <taybart@email.com>",
		About:   "check for defined env vars in a project or file",
		Args: map[string]*args.Arg{
			"files": {
				Short: "f",
				Long:  "files",
				Help:  "Comma seperated files to check (./main.go,./util.go)",
			},
			"directory": {
				Short: "d",
				Long:  "directory",
				Help:  "Scan directory",
			},
			"validate": {
				Short: "v",
				Long:  "validate",
				Help:  "File to validate env config against",
			},
			"print": {
				Short:   "p",
				Long:    "print",
				Help:    "Print contents in env file format, will add file tags above each env",
				Default: false,
			},
		},
	}
)

func Scan(app args.App) error {
	files := []string{}

	switch {
	case isPiped(): // piped output into the cli
		reader := bufio.NewReader(os.Stdin)
		var fns []rune
		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			fns = append(fns, input)
		}
		files = strings.Split(string(fns), "\n")

	case app.Get("directory").IsSet(): // directory specified
		re := regexp.MustCompile(`[[:alnum:]\/\._\-]+.go$`)
		dir := app.String("directory")
		err := filepath.Walk(dir, func(path string, _ os.FileInfo, e error) error {
			if e != nil {
				return e
			}
			if re.Match([]byte(path)) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return err
		}

	case app.Get("files").IsSet(): // csv of files
		files = strings.Split(app.String("files"), ",")
	default:
		return errors.New("no files specified")
	}

	// Get down to buisness
	v := newVisitor()
	for _, f := range files {
		node, err := v.Load(f)
		if err != nil {
			continue
		}
		ast.Inspect(node, v.Visit)
	}
	if app.Get("validate").IsSet() {
		log.Debug("Should Validate", app.String("validate"))
		foundEnv, optional := v.EnvToMap()

		envToTest, err := parseEnvFile(app.String("validate"))
		if err != nil {
			return err
		}
		missing := []string{}
		usingDefault := []string{}
		for k, v := range foundEnv {
			if _, ok := envToTest[k]; !ok {
				if optional[k] {
					continue
				}
				if v != "" {
					usingDefault = append(usingDefault, fmt.Sprintf("%s=\"%s\"", k, v))
					continue
				}
				missing = append(missing, k)
			}
		}
		sort.Strings(missing)
		sort.Strings(usingDefault)

		if len(missing) > 0 {
			m := ""
			for _, v := range missing {
				m += fmt.Sprintf("%s%s%s\n", log.Red, v, log.Reset)
			}
			log.Errorf("Missing required env\n%s", m)
		}

		if len(usingDefault) > 0 {
			for _, v := range usingDefault {
				k, d := env.GetDefault(v)
				log.Warnf("Using default value for %s=%s\n", k, strings.Trim(d, `"`))
			}
		}
		return nil

	}

	if app.Get("print").IsSet() {
		fmt.Println(v.EnvByFile())
		return nil
	}
	fmt.Println(v.ToEnvFile())
	return nil
}

// Check if program has data piped to it
func isPiped() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	return info.Mode()&os.ModeCharDevice == 0
}
