package main

import (
	"bufio"
	"errors"
	"go/ast"
	"log"
	"os"
	"regexp"
)

var (
	envRE = regexp.MustCompile(`([[:word:]]+)="?(.*)?"?`)
)

// isIdent: Checks that expr is an idenifier
func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// isPkgDot: Checks that call belongs to pkg
func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

// getStringArray: Convert ast.CompositeLiteral to string array
func getStringArray(l *ast.CompositeLit) ([]string, error) {
	arr := []string{}
	t, ok := l.Type.(*ast.ArrayType)
	if !ok {
		return arr, errors.New("not an array")
	}
	if s, ok := t.Elt.(*ast.Ident); !ok || s.Name != "string" {
		return arr, errors.New("array of incorrect type")
	}

	for _, i := range l.Elts {
		if v, ok := i.(*ast.BasicLit); ok {
			arr = append(arr, v.Value)
		}
	}

	return arr, nil
}

// dedupe: Deduplicate array
func dedupe(in []string) []string {
	seen := make(map[string]bool)
	dd := []string{}
	for _, i := range in {
		if !seen[i] {
			seen[i] = true
			dd = append(dd, i)
		}
	}
	return dd

}

func parseEnvFile(filename string) (envFile map[string]string, err error) {
	envFile = make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res := envRE.FindAllStringSubmatch(scanner.Text(), -1)
		if len(res) > 0 {
			envFile[res[0][1]] = res[0][2]
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}
