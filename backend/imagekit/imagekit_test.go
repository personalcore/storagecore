package imagekit

import (
	"testing"

	"github.com/personalcore/storagecore/fstest"
	"github.com/personalcore/storagecore/fstest/fstests"
)

func TestIntegration(t *testing.T) {
	debug := true
	fstest.Verbose = &debug
	fstests.Run(t, &fstests.Opt{
		RemoteName:      "TestImageKit:",
		NilObject:       (*Object)(nil),
		SkipFsCheckWrap: true,
	})
}
