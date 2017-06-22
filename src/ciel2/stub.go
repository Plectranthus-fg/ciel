package main

import (
	"ciel-driver"
	"errors"
	"os"
	"os/exec"
)

func updateStub(c *ciel.Container) error {
	c.Unmount()
	os.RemoveAll(c.Fs.UpperDir)
	defer func() {
		c.Fs.MaskCache = false
		c.Fs.MaskBuildKit = false
		c.Fs.MaskStubConfig = false
		c.Fs.MaskStub = false
		c.Unmount()
	}()

	c.Fs.MaskCache = true
	c.Fs.MaskBuildKit = true
	c.Fs.MaskStubConfig = false
	c.Fs.MaskStub = false
	if ec := c.Command("apt update -y"); ec != 0 {
		return errors.New("apt update: failed")
	}
	if ec := c.Command("apt full-upgrade -y"); ec != 0 {
		return errors.New("apt full-upgrade: failed")
	}

	c.Fs.MaskCache = true
	c.Fs.MaskBuildKit = true
	c.Fs.MaskStubConfig = true
	c.Fs.MaskStub = false
	c.Unmount()
	if err := cleanStub(c); err != nil {
		return err
	}
	return mergeStub(c)
}

func cleanStub(c *ciel.Container) error {
	return clean(c, []string{
		`^/dev`,
		`^/efi`,
		`^/etc`,
		`^/run`,
		`^/usr`,
		`^/var/lib/dpkg`,
		`^/var/log/journal$`,
		`^/root`,
		`^/home/aosc`,
		`/\.updated$`,
	}, []string{
		`^/etc/.*-$`,
		`^/etc/ssh/ssh_host_.*`,
		`^/var/lib/dpkg/status-old`,
	}, func(path string, info os.FileInfo, err error) error {
		if err := os.RemoveAll(path); err != nil {
			println(path, err.Error())
		}
		return nil
	})
}

func mergeStub(c *ciel.Container) error {
	tmp := c.Fs.Stub + "." + randomFilename()
	c.Mount()
	{
		os.Mkdir(tmp, 0775)
		cmd := exec.Command("/bin/cp", "-prT", c.Fs.Target, tmp)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	c.Unmount()

	if err := os.RemoveAll(c.Fs.Stub); err != nil {
		return err
	}
	if err := os.Rename(tmp, c.Fs.Stub); err != nil {
		return err
	}
	return nil
}