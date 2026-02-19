// Package cmdtest creates a testable interface to rclone main
//
// The interface is used to perform end-to-end test of
// commands, flags, environment variables etc.
package cmdtest

// The rest of this file is a 1:1 copy from rclone.go

import (
	_ "github.com/personalcore/storagecore/backend/all" // import all backends
	"github.com/personalcore/storagecore/cmd"
	_ "github.com/personalcore/storagecore/cmd/all"    // import all commands
	_ "github.com/personalcore/storagecore/lib/plugin" // import plugins
)

func main() {
	cmd.Main()
}
