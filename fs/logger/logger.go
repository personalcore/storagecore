// Package logger implements testing for the sync (and bisync) logger
package logger

import (
	_ "github.com/personalcore/storagecore/backend/all" // import all backends
	"github.com/personalcore/storagecore/cmd"
	_ "github.com/personalcore/storagecore/cmd/all"    // import all commands
	_ "github.com/personalcore/storagecore/lib/plugin" // import plugins
)

// Main enables the testscript package. See:
// https://bitfieldconsulting.com/golang/cli-testing
// https://pkg.go.dev/github.com/rogpeppe/go-internal@v1.11.0/testscript
func Main() {
	cmd.Main()
}
