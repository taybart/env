package scan

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/taybart/env"
)

type visitor struct {
	decls       map[string][]string
	env         map[string][]string
	fset        *token.FileSet
	packageName string
	fn          string
}

func newVisitor() visitor {
	// generate tokens
	fset := token.NewFileSet()
	return visitor{
		decls: make(map[string][]string),
		env:   make(map[string][]string),
		fset:  fset,
	}
}

// Load: Loads file node and checks that the env package is imported
func (v *visitor) Load(filename string) (ast.Node, error) {
	v.fn = filename
	// get node to parse
	node, err := parser.ParseFile(v.fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// check if the package is actually imported
	usesEnv := false
	for _, i := range node.Imports {
		if strings.Contains(i.Path.Value, "github.com/taybart/env") {
			usesEnv = true
			v.packageName = "env"
			if i.Name != nil {
				v.packageName = i.Name.String()
			}
		}
	}
	if !usesEnv {
		return nil, errors.New("does not use env package")
	}
	return node, nil
}

// Visit: Implements ast.Visit to be passed to ast.Inspect
func (v *visitor) Visit(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.AssignStmt: // Line contains an assignment with :=
		for _, name := range n.Lhs {
			if ident, ok := name.(*ast.Ident); ok {
				if ident.Name == "_" {
					return true
				}
				if ident.Obj != nil && ident.Obj.Pos() == ident.Pos() {
					if compLit, ok := n.Rhs[0].(*ast.CompositeLit); ok {
						arr, err := getStringArray(compLit)
						if err != nil {
							return true
						}
						v.decls[ident.Name] = arr
					}
				}
			}
		}
	case *ast.GenDecl: // line contains a declaration
		if n.Tok == token.CONST || n.Tok == token.VAR { // we only care about var/const
			for _, spec := range n.Specs {
				switch spec := spec.(type) {
				case *ast.ValueSpec:
					for _, name := range spec.Names {
						if vspec, ok := name.Obj.Decl.(*ast.ValueSpec); ok { // get the
							if len(vspec.Values) > 0 {
								if cl, ok := vspec.Values[0].(*ast.CompositeLit); ok {
									arr, err := getStringArray(cl)
									if err != nil {
										return true
									}
									v.decls[name.Name] = arr
								}
							}
						}
					}
				}
			}
		}
	case *ast.CallExpr: // line contains a function call
		if isPkgDot(n.Fun, v.packageName, "Set") { // check that function call is env.Set() (handles import renames)
			switch arg := n.Args[0].(type) {
			case *ast.CompositeLit: // function was passed anon []string
				arr, _ := getStringArray(arg)
				v.env[v.fn] = append(v.env[v.fn], arr...)
			case *ast.Ident: // function was passed a variable
				v.env[v.fn] = append(v.env[v.fn], v.decls[arg.Name]...)
			}
		}
	}
	return true
}

func (v visitor) EnvToMap() (map[string]string, map[string]bool) {
	e := []string{}
	for _, ev := range v.env {
		e = append(e, ev...)
	}
	e = dedupe(e)

	optional := env.GetOptional(e)

	envmap := make(map[string]string)
	for _, k := range e {
		key, val := env.GetDefault(k[1 : len(k)-1])
		envmap[key] = val
	}
	return envmap, optional
}

func (v visitor) ToEnvFile() string {
	e := []string{}
	for _, en := range v.env {
		e = append(e, en...)
	}
	e = dedupe(e)

	output := ""
	optional := env.GetOptional(e)
	for i, k := range e {
		key, val := env.GetDefault(k[1 : len(k)-1])
		if optional[key] {
			val = "Value is marked as optional"
		}
		output += fmt.Sprintf("%s=\"%s\"", key, val)
		if i < len(e)-1 {
			output += "\n"
		}
	}
	return output
}

func (v visitor) EnvByFile() string {
	output := ""
	for ns, e := range v.env {
		optional := env.GetOptional(v.env[ns])
		output += fmt.Sprintf("#%s\n", ns)
		for _, k := range e {
			key, val := env.GetDefault(k[1 : len(k)-1])
			if optional[key] {
				val = "Value is marked as optional"
			}
			output += fmt.Sprintf("%s=\"%s\"\n", key, val)
		}
		output += "\n"
	}
	return output
}
