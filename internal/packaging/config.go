package packaging

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"ciel/display"
	"ciel/internal/abstract"
)

const (
	DefaultEditor = "/usr/bin/editor"
)

func EditSourceList(global bool, i abstract.Instance, c abstract.Container) {
	var root string
	if global {
		root = c.DistDir()
	} else {
		root = i.MountPoint()
	}
	editor := editor()
	cmd := exec.Command(editor, path.Join(root, "/etc/apt/sources.list"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func SetTreePath(i abstract.Instance, tree string) {
	root := i.MountPoint()
	config := `[default]` + "\n"
	config += `location = ` + path.Clean(tree) + "\n"
	d.ITEM("set tree path")
	err := ioutil.WriteFile(path.Join(root, "/etc/acbs/forest.conf"), []byte(config), 0644)
	d.ERR(err)
}

func DisableDNSSEC(i abstract.Instance) {
	root := i.MountPoint()
	config := `[Resolve]` + "\n"
	config += `DNSSEC=no` + "\n"
	d.ITEM("disable DNSSEC")
	err := ioutil.WriteFile(path.Join(root, "/etc/systemd/resolved.conf"), []byte(config), 0644)
	d.ERR(err)
}

func SetMaintainer(i abstract.Instance, person string) {
	root := i.MountPoint()
	config := `#!/bin/bash` + "\n"
	config += `ABMPM=dpkg` + "\n"
	config += `ABAPMS=` + "\n"
	config += `MTER="` + person + `"` + "\n"
	config += `ABINSTALL=` + "\n"
	d.ITEM("set maintainer")
	err := ioutil.WriteFile(path.Join(root, "/usr/lib/autobuild3/etc/autobuild/ab3cfg.sh"), []byte(config), 0644)
	d.ERR(err)
}

func editor() string {
	if s := os.Getenv("VISUAL"); s != "" {
		return s
	}
	if s := os.Getenv("EDITOR"); s != "" {
		return s
	}
	return DefaultEditor
}
