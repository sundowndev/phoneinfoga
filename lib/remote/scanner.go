package remote

import (
	"fmt"
	"os"
	"plugin"

	"github.com/sundowndev/phoneinfoga/v2/lib/number"
)

type ScannerOptions map[string]interface{}

func (o ScannerOptions) GetStringEnv(k string) string {
	if v, ok := o[k].(string); ok {
		return v
	}
	return os.Getenv(k)
}

type Plugin interface {
	Lookup(string) (plugin.Symbol, error)
}

type Scanner interface {
	Name() string
	Description() string
	DryRun(number.Number, ScannerOptions) error
	Run(number.Number, ScannerOptions) (interface{}, error)
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
