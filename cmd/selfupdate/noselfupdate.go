//go:build noselfupdate

package selfupdate

import (
	"github.com/personalcore/storagecore/lib/buildinfo"
)

func init() {
	buildinfo.Tags = append(buildinfo.Tags, "noselfupdate")
}
