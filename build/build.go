package build

import (
	"fmt"
	"os"
)

// Version is the corresponding release tag
var Version = "dev"

// Commit is the corresponding Git commit
var Commit = "dev"

func IsRelease() bool {
	return String() != "dev-dev"
}

func String() string {
	return fmt.Sprintf("%s-%s", Version, Commit)
}

func IsDemo() bool {
	return os.Getenv("PHONEINFOGA_DEMO") == "true"
}
