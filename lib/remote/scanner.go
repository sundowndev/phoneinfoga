package remote

import (
	"fmt"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"os"
	"plugin"
)

type Plugin interface {
	Lookup(string) (plugin.Symbol, error)
}

type Scanner interface {
	Scan(*number.Number) (interface{}, error)
	ShouldRun() bool
	Identifier() string
}

func parseEntryFunc(p Plugin) (Scanner, error) {
	symbol, err := p.Lookup("NewScanner")
	if err != nil {
		return nil, fmt.Errorf("exported function NewScanner not found")
	}

	fn, ok := symbol.(func() Scanner)
	if !ok {
		return nil, fmt.Errorf("exported function NewScanner does not follow the remote.Scanner interface")
	}
	return fn(), nil
}

func OpenPlugin(path string) (Scanner, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("given path %s does not exist", path)
	}

	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("given plugin %s is not valid", path)
	}

	return parseEntryFunc(p)
}
