package remote

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

type ScannerResult map[string]interface{}

type Scanner interface {
	Scan(*number.Number) (ScannerResult, error)
	ShouldRun() bool
	Identifier() string
}

func InitScanners(remote *Library) {
	numverifySupplier := suppliers.NewNumverifySupplier()

	remote.AddScanner(NewNumverifyScanner(numverifySupplier))
}
