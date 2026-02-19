//go:build cmount && ((linux && cgo) || (darwin && cgo) || (freebsd && cgo) || (openbsd && cgo) || windows) && (!race || !windows)

// Package cmount implements a FUSE mounting system for rclone remotes.
//
// FIXME this doesn't work with the race detector under Windows either
// hanging or producing lots of differences.

package cmount

import (
	"runtime"
	"testing"

	"github.com/personalcore/storagecore/fstest/testy"
	"github.com/personalcore/storagecore/vfs/vfscommon"
	"github.com/personalcore/storagecore/vfs/vfstest"
)

func TestMount(t *testing.T) {
	// Disable tests under macOS and the CI since they are locking up
	if runtime.GOOS == "darwin" {
		testy.SkipUnreliable(t)
	}
	vfstest.RunTests(t, false, vfscommon.CacheModeOff, true, mount)
}
