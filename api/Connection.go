package api

import (
	"os"

	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (i *InstHandler) IsConnected() bool {
	if i.LibvirtInst == nil {
		return false
	}
	alive, err := i.LibvirtInst.IsAlive()
	// As IsAlive do a ping to libvirt daemon
	// better implement a more efficient way to check the connection status of libvirt daemon
	// e.g. heartbeat
	return err == nil && alive
}

func (i *InstHandler) LibvirtConnection() {
	libvirtInst, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		i.Logger.Panic("innitial connection for libvirt daemon failed", zap.Int("pid", os.Getegid()))
	}
	i.LibvirtInst = libvirtInst
	i.Logger.Info("Libvirt Coonnection succefully done.", zap.Int("pid", os.Getegid()))

	defer i.Logger.Sync()
}
