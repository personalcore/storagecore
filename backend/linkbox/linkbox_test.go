// Test Linkbox filesystem interface
package linkbox_test

import (
	"testing"

	"github.com/personalcore/storagecore/backend/linkbox"
	"github.com/personalcore/storagecore/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestLinkbox:",
		NilObject:  (*linkbox.Object)(nil),
		// Linkbox doesn't support leading dots for files
		SkipLeadingDot: true,
	})
}
