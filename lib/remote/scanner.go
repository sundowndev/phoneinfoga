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
	Name() string
	Scan(number.Number) (interface{}, error)
	ShouldRun(number.Number) bool
}

func OpenPlugin(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("given path %s does not exist", path)
	}

	_, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("given plugin %s is not valid", path)
	}

	return nil
}
