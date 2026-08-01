package remote

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

func InitScanners(remote *Library) {
	numverifySupplier := suppliers.NewNumverifySupplier()
	ovhSupplier := suppliers.NewOVHSupplier()
	nomoroboSupplier := suppliers.NewNomoroboSupplier()
	shouldianswerSupplier := suppliers.NewShouldIAnswerSupplier()

	remote.AddScanner(NewLocalScanner())
	remote.AddScanner(NewNumverifyScanner(numverifySupplier))
	remote.AddScanner(NewGoogleSearchScanner())
	remote.AddScanner(NewOVHScanner(ovhSupplier))
	remote.AddScanner(NewGoogleCSEScanner(nil))
	remote.AddScanner(NewNomoroboScanner(nomoroboSupplier))
	remote.AddScanner(NewShouldIAnswerScanner(shouldianswerSupplier))

	remote.LoadPlugins()
}
