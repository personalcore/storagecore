// Test Seafile filesystem interface
package seafile_test

import (
	"testing"

	"github.com/personalcore/storagecore/backend/seafile"
	"github.com/personalcore/storagecore/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSeafile:",
		NilObject:  (*seafile.Object)(nil),
	})
}
