package scan

import (
	"fmt"
	"slices"

	"github.com/taybart/env"
)

type EnvVar struct {
	Value      string
	Optional   bool
	HasDefault bool
}
type Env struct {
	Values map[string]EnvVar
	v      *visitor
}

func NewEnv() Env {
	return Env{
		Values: make(map[string]EnvVar),
	}
}

func (e Env) Equal(cmp Env) bool {
	if len(e.Values) != len(cmp.Values) {
		fmt.Println("env lengths not equal")
		return false
	}
	for k, v := range e.Values {
		if v.Value != cmp.Values[k].Value ||
			v.Optional != cmp.Values[k].Optional ||
			v.HasDefault != cmp.Values[k].HasDefault {
			fmt.Println(k, "not equal")
			return false
		}
	}
	return true
}
func (e Env) ToFile() string {
	output := ""

	// force alphabetical order of map
	order := []string{}
	for k := range e.Values {
		order = append(order, k)
	}
	slices.Sort(order)

	for i, v := range order {
		entry := e.Values[v]
		val := entry.Value
		if entry.Optional {
			val = "value is marked as optional"
		}
		output += fmt.Sprintf("%s=\"%s\"", v, val)
		if i < len(order)-1 {
			output += "\n"
		}
	}
	return output
}

func (e Env) EnvByFile() string {
	if e.v == nil {
		return ""
	}
	output := ""
	for ns, en := range e.v.env {
		optional := env.GetOptional(e.v.env[ns])
		output += fmt.Sprintf("#%s\n", ns)
		for _, k := range en {
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
