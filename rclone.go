// Sync files and directories to and from local and remote object stores
//
// Nick Craig-Wood <nick@craig-wood.com>
package main

import (
	_ "github.com/personalcore/storagecore/backend/all" // import all backends
	"github.com/personalcore/storagecore/cmd"
	_ "github.com/personalcore/storagecore/cmd/all"    // import all commands
	_ "github.com/personalcore/storagecore/lib/plugin" // import plugins
)

func main() {
	cmd.Main()
}
