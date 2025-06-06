package scan

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/taybart/env"
	"github.com/taybart/log"
)

type Config struct {
	Dir        string `arg:"directory"`
	Files      string `arg:"files"`
	Tags       string `arg:"tags"`
	PrintFiles bool   `arg:"print"`
	Validate   string `arg:"validate"`
}

// TODO: output map, warn that project is not go project, expand_path
func Scan(config Config) (Env, error) {
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

	case config.Dir != "": // directory specified
		re := regexp.MustCompile(`[[:alnum:]\/\._\-]+.go$`)
		err := filepath.Walk(config.Dir, func(path string, _ os.FileInfo, e error) error {
			if e != nil {
				return e
			}
			if re.Match([]byte(path)) {
				// check if build tags apply
				if config.Tags != "" {
					ok, err := checkBuildTags(strings.Split(config.Tags, ","), path)
					if err != nil {
						return err
					}
					if !ok {
						return nil
					}
				}
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return Env{}, err
		}

	case config.Files != "": // csv of files
		files = strings.Split(config.Files, ",")
	default:
		return Env{}, errors.New("no files specified")
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
	if config.Validate != "" {
		log.Debug("Should Validate", config.Validate)
		foundEnv, optional := v.EnvToMap()

		envToTest, err := parseEnvFile(config.Validate)
		if err != nil {
			return Env{}, err
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
		return Env{}, nil

	}

	return v.Finish(), nil

	// if config.PrintFiles {
	// 	return v.EnvByFile(), nil
	// }
	// return v.ToEnvFile(), nil
}

// Check if program has data piped to it
func isPiped() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	return info.Mode()&os.ModeCharDevice == 0
}

func checkBuildTags(tags []string, path string) (bool, error) {
	context := build.Default
	context.BuildTags = tags
	return context.MatchFile(filepath.Dir(path), filepath.Base(path))
}
