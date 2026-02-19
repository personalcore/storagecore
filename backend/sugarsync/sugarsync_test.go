// Test Sugarsync filesystem interface
package sugarsync_test

import (
	"testing"

	"github.com/personalcore/storagecore/backend/sugarsync"
	"github.com/personalcore/storagecore/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSugarSync:Test",
		NilObject:  (*sugarsync.Object)(nil),
	})
}
