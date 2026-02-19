//go:build !plan9 && !solaris

package iclouddrive_test

import (
	"testing"

	"github.com/personalcore/storagecore/backend/iclouddrive"
	"github.com/personalcore/storagecore/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestICloudDrive:",
		NilObject:  (*iclouddrive.Object)(nil),
	})
}
