package scripts

import (
	"fmt"
	"sort"
)

func (sh *ScriptHandler) GetDocstring(function string) string {
	return sh.RegisteredFunctions[function].Docstring
}
func (sh *ScriptHandler) MakeDocs() {
	keys := sh.FunctionNames()

	output := ""
	for _, k := range keys {
		output += fmt.Sprintf("%s: %s", k, sh.RegisteredFunctions[k].Docstring)
	}
}

func (sh *ScriptHandler) FunctionNames() []string {
	keys := make([]string, len(sh.RegisteredFunctions))
	i := 0
	for k, _ := range sh.RegisteredFunctions {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	return keys
}
