//go:build unix

// The serving is tested in cmd/nfsmount - here we test anything else
package nfs

import (
	"testing"

	_ "github.com/personalcore/storagecore/backend/local"
	"github.com/personalcore/storagecore/cmd/serve/servetest"
	"github.com/personalcore/storagecore/fs/rc"
)

func TestRc(t *testing.T) {
	servetest.TestRc(t, rc.Params{
		"type":           "nfs",
		"vfs_cache_mode": "off",
	})
}
