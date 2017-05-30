package ci

import (
	"encoding/base64"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type ContainerFilesystem struct {
	Base string // overlayfs: lowerdir 1

	overlay string // overlayfs: lowerdir 2
	diff    string // overlayfs: upperdir
	work    string // overlayfs: workdir

	TargetDir string // overlayfs: target
}

func InitFilesystem(bkdir string) *ContainerFilesystem {
	fs := &ContainerFilesystem{Base: bkdir}
	fs.overlay = fs.Base + ".overlay"
	fs.diff = fs.Base + ".diff"
	fs.work = fs.Base + ".temp"
	os.Mkdir(fs.diff, 0755)
	os.Mkdir(fs.work, 0755)
	return fs
}

func (fs *ContainerFilesystem) Startup() error {
	rd := make([]byte, 8)
	if _, err := rand.Read(rd); err != nil {
		return err
	}
	fs.TargetDir = os.TempDir() + "/ciel." + base64.RawURLEncoding.EncodeToString(rd)
	os.Mkdir(fs.TargetDir, 0755)
	if _, err := os.Stat(fs.overlay); os.IsNotExist(err) {
		return mount(fs.TargetDir, fs.diff, fs.work, fs.Base)
	} else {
		return mount(fs.TargetDir, fs.diff, fs.work, fs.overlay, fs.Base)
	}
}

func (fs *ContainerFilesystem) Shutdown() error {
	if err := unmount(fs.TargetDir); err != nil {
		return err
	}
	if err := os.Remove(fs.TargetDir); err != nil {
		return err
	}
	if err := os.RemoveAll(fs.work); err != nil {
		return err
	}
	return nil
}

func mount(path string, upperdir string, workdir string, lowerdirs ...string) error {
	return syscall.Mount("overlay", path, "overlay", 0,
		"lowerdir="+strings.Join(lowerdirs, ":")+",upperdir="+upperdir+",workdir="+workdir)
}

func unmount(path string) error {
	return syscall.Unmount(path, 0)
}

func (fs *ContainerFilesystem) DiffDir(path string) string {
	return fs.diff + path
}

func (fs *ContainerFilesystem) OverlayDir(path string) string {
	return fs.overlay + path
}

func (fs *ContainerFilesystem) Merge(path string) error {
	os.Mkdir(fs.overlay, 0755)
	err := filepath.Walk(fs.DiffDir(path), func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		rel, err := filepath.Rel(fs.DiffDir("/"), path)
		if err != nil {
			return err
		}
		rel = "/" + rel
		if info.IsDir() {
			os.MkdirAll(fs.OverlayDir(rel), 755)
		}
		if err := os.Rename(path, fs.OverlayDir(rel)); err == nil {
			log.Println("clean: merge", rel)
		}
		return nil
	})
	return err
}
