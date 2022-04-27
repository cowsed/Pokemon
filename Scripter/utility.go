package scripts

import (
	"fmt"
	"sort"
)

func (sh *ScriptEngine) GetDocstring(function string) string {
	return sh.RegisteredFunctions[function].Docstring
}
func (sh *ScriptEngine) MakeDocs() {
	keys := sh.FunctionNames()

	output := ""
	for _, k := range keys {
		output += fmt.Sprintf("%s: %s", k, sh.RegisteredFunctions[k].Docstring)
	}
}

func (sh *ScriptEngine) FunctionNames() []string {
	keys := make([]string, len(sh.RegisteredFunctions))
	i := 0
	for k := range sh.RegisteredFunctions {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	return keys
}
