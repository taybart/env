package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	filesFlag    string
	dirFlag      string
	prettyPrint  bool
	validateFlag string
)

func init() {
	flag.StringVar(&filesFlag, "f", "", "csv of files to check ex. -f ./main.go,./util.go")
	flag.StringVar(&dirFlag, "d", "", "directory containing the project")
	flag.StringVar(&validateFlag, "validate", "", "file to validate env against")
	flag.BoolVar(&prettyPrint, "p", false, "print env by file")
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func run() error {
	flag.Parse()

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

	case dirFlag != "": // directory specified
		re := regexp.MustCompile(`[[:alnum:]\/\._\-]+.go$`)
		err := filepath.Walk(dirFlag, func(path string, info os.FileInfo, e error) error {
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

	case filesFlag != "": // csv of files
		files = strings.Split(filesFlag, ",")
	default:
		return errors.New("No files specified")
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

	if validateFlag != "" {
		foundEnv, optional := v.EnvToMap()

		envToTest, err := parseEnvFile(validateFlag)
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
			fmt.Println("~~Missing required~~")
			for _, v := range missing {
				fmt.Println("\t", v)
			}
		}

		if len(usingDefault) > 0 {
			fmt.Println("\n~~Using defaults~~")
			for _, v := range usingDefault {
				fmt.Println("\t", v)
			}
		}
		return nil

	}

	if prettyPrint {
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
