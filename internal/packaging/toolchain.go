package packaging

import (
	"os"
	"path"

	"ciel/display"
	"ciel/internal/abstract"
)

const (
	AB3Path  = "/usr/bin/autobuild"
	ACBSPath = "/usr/bin/acbs-build"
)

type ToolChain struct {
	AB   bool
	ACBS bool
}

func DetectToolChain(i abstract.Instance) *ToolChain {
	root := i.MountPoint()
	tc := &ToolChain{}
	d.ITEM("detect autobuild3")
	tc.AB = exists(root, AB3Path)
	d.ITEM("detect acbs")
	tc.ACBS = exists(root, ACBSPath)
	return tc
}

func exists(root, target string) bool {
	_, err := os.Stat(path.Join(root, target))
	if os.IsNotExist(err) {
		d.FAILED()
		return false
	} else if err == nil {
		d.OK()
		return true
	} else {
		d.FAILED_BECAUSE(err.Error())
		return false
	}
}
