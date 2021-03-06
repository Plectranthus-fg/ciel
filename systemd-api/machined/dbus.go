package machined

import (
	"ciel/systemd-api"
	"github.com/godbus/dbus"
)

const Dest = "org.freedesktop.machine1"

func Object(path dbus.ObjectPath) dbus.BusObject {
	return systemd.Conn.Object(Dest, path)
}
