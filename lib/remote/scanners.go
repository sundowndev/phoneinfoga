package remote

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

type Scanner interface {
	Scan(*number.Number) (interface{}, error)
	ShouldRun() bool
	Identifier() string
}

func InitScanners(remote *Library) {
	numverifySupplier := suppliers.NewNumverifySupplier()
	ovhSupplier := suppliers.NewOVHSupplier()

	remote.AddScanner(NewLocalScanner())
	remote.AddScanner(NewNumverifyScanner(numverifySupplier))
	remote.AddScanner(NewGoogleSearchScanner())
	remote.AddScanner(NewOVHScanner(ovhSupplier))
}
