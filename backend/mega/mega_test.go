// Test Mega filesystem interface
package mega_test

import (
	"testing"

	"github.com/personalcore/storagecore/backend/mega"
	"github.com/personalcore/storagecore/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestMega:",
		NilObject:  (*mega.Object)(nil),
	})
}
