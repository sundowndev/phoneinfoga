package remote

import (
	"fmt"
	"os"
	"plugin"

	"github.com/sundowndev/phoneinfoga/v2/lib/number"
)

type Plugin interface {
	Lookup(string) (plugin.Symbol, error)
}

type Scanner interface {
	Name() string
	Description() string
	DryRun(number.Number) error
	Run(number.Number) (interface{}, error)
}

func OpenPlugin(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("given path %s does not exist", path)
	}

	_, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("given plugin %s is not valid: %v", path, err)
	}

	return nil
}
