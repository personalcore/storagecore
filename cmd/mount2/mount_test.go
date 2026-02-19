//go:build linux

package mount2

import (
	"testing"

	"github.com/personalcore/storagecore/vfs/vfscommon"
	"github.com/personalcore/storagecore/vfs/vfstest"
)

func TestMount(t *testing.T) {
	vfstest.RunTests(t, false, vfscommon.CacheModeOff, true, mount)
}
