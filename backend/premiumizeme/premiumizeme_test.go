// Test filesystem interface
package premiumizeme_test

import (
	"testing"

	"github.com/personalcore/storagecore/backend/premiumizeme"
	"github.com/personalcore/storagecore/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestPremiumizeMe:",
		NilObject:  (*premiumizeme.Object)(nil),
	})
}
