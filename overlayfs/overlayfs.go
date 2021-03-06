package overlayfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"ciel/display"
)

type Instance struct {
	MountPoint string
	Layers     []string
}

const TmpDirSuffix = ".tmp"

// https://www.kernel.org/doc/Documentation/filesystems/overlayfs.txt
//
// Multiple lower layers
// ---------------------
//
// Multiple lower layers can now be given using the the colon (":") as a
// separator character between the directory names.  For example:
//
//   mount -t overlay overlay -olowerdir=/lower1:/lower2:/lower3 /merged
//
// As the example shows, "upperdir=" and "workdir=" may be omitted.  In
// that case the overlay will be read-only.
//
// The specified lower directories will be stacked beginning from the
// rightmost one and going left.  In the above example lower1 will be the
// top, lower2 the middle and lower3 the bottom layer.

func (i *Instance) MountLocal() error {
	var localInst Instance
	localInst = *i
	localInst.Layers = localInst.Layers[:len(localInst.Layers)-1]
	return localInst.Mount(false)
}

func (i *Instance) Mount(readOnly bool) error {
	var option string
	var layers = make([]string, len(i.Layers))
	for index := range layers {
		layers[index] = filepath.Clean(i.Layers[len(layers)-1-index]) // reverse i.Layers and assign it to layers
	}
	if readOnly {
		option = "lowerdir=" + strings.Join(layers, ":")
	} else {
		olfsLowerdirs := layers[1:]
		olfsUpperdir := layers[0]
		olfsWorkdir := olfsUpperdir + TmpDirSuffix
		os.MkdirAll(olfsWorkdir, 0755)
		option =
			"lowerdir=" + strings.Join(olfsLowerdirs, ":") +
				",upperdir=" + olfsUpperdir +
				",workdir=" + olfsWorkdir
	}
	os.MkdirAll(i.MountPoint, 0755)
	err := syscall.Mount("overlay", i.MountPoint, "overlay", 0, option)
	return err
}

func (i *Instance) Unmount() error {
	err := syscall.Unmount(i.MountPoint, 0)
	if err == nil {
		if len(i.Layers) > 0 {
			os.RemoveAll(filepath.Clean(i.Layers[len(i.Layers)-2]) + TmpDirSuffix)
			os.RemoveAll(filepath.Clean(i.Layers[len(i.Layers)-1]) + TmpDirSuffix)
		}
	}
	return err
}

func (i *Instance) Rollback() error {
	d.ITEM("get diff dir")
	layers := i.Layers
	dir := layers[len(layers)-1]
	d.Println(d.C(d.WHITE, dir))

	d.ITEM("clean diff dir")
	fi, err := ioutil.ReadDir(dir)
	if err == nil {
		for _, f := range fi {
			err = os.RemoveAll(filepath.Join(dir, f.Name()))
		}
		d.ERR(err)
	} else if os.IsNotExist(err) {
		d.SKIPPED()
	} else {
		d.FAILED_BECAUSE(err.Error())
		return err
	}
	return nil
}
