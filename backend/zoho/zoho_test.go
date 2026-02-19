// Test Zoho filesystem interface
package zoho_test

import (
	"testing"

	"github.com/personalcore/storagecore/backend/zoho"
	"github.com/personalcore/storagecore/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName:      "TestZoho:",
		SkipInvalidUTF8: true,
		NilObject:       (*zoho.Object)(nil),
	})
}
