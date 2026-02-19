package filen

import (
	"testing"

	"github.com/personalcore/storagecore/fstest/fstests"
)

func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestFilen:",
		NilObject:  (*Object)(nil),
	})
}
