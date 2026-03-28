package virtxml

import (
	"testing"

	libvirtxml "libvirt.org/libvirt-go-xml"
)

func emptyDeviceList() *libvirtxml.DomainDeviceList {
	return &libvirtxml.DomainDeviceList{}
}

func TestDisk_ApplyTo(t *testing.T) {
	devs := emptyDeviceList()
	d := Disk{libvirtxml.DomainDisk{Device: "disk"}}
	d.ApplyTo(devs)

	if len(devs.Disks) != 1 {
		t.Fatalf("expected 1 disk, got %d", len(devs.Disks))
	}
	if devs.Disks[0].Device != "disk" {
		t.Errorf("device: expected disk, got %s", devs.Disks[0].Device)
	}
}

func TestDisk_ApplyTo_Multiple(t *testing.T) {
	devs := emptyDeviceList()
	Disk{libvirtxml.DomainDisk{Device: "disk"}}.ApplyTo(devs)
	Disk{libvirtxml.DomainDisk{Device: "cdrom"}}.ApplyTo(devs)

	if len(devs.Disks) != 2 {
		t.Fatalf("expected 2 disks, got %d", len(devs.Disks))
	}
}

func TestNetworkInterface_ApplyTo(t *testing.T) {
	devs := emptyDeviceList()
	iface := NetworkInterface{libvirtxml.DomainInterface{
		MAC:   &libvirtxml.DomainInterfaceMAC{Address: "52:54:00:ab:cd:ef"},
		Model: &libvirtxml.DomainInterfaceModel{Type: "virtio"},
	}}
	iface.ApplyTo(devs)

	if len(devs.Interfaces) != 1 {
		t.Fatalf("expected 1 interface, got %d", len(devs.Interfaces))
	}
	if devs.Interfaces[0].MAC.Address != "52:54:00:ab:cd:ef" {
		t.Errorf("mac: expected 52:54:00:ab:cd:ef, got %s", devs.Interfaces[0].MAC.Address)
	}
}

func TestChannel_ApplyTo(t *testing.T) {
	devs := emptyDeviceList()
	ch := Channel{libvirtxml.DomainChannel{
		Target: &libvirtxml.DomainChannelTarget{
			VirtIO: &libvirtxml.DomainChannelTargetVirtIO{Name: "org.qemu.guest_agent.0"},
		},
	}}
	ch.ApplyTo(devs)

	if len(devs.Channels) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(devs.Channels))
	}
	if devs.Channels[0].Target.VirtIO.Name != "org.qemu.guest_agent.0" {
		t.Errorf("channel name: expected org.qemu.guest_agent.0, got %s", devs.Channels[0].Target.VirtIO.Name)
	}
}

func TestSerial_ApplyTo(t *testing.T) {
	devs := emptyDeviceList()
	port := uint(0)
	Serial{libvirtxml.DomainSerial{
		Source: &libvirtxml.DomainChardevSource{Pty: &libvirtxml.DomainChardevSourcePty{}},
		Target: &libvirtxml.DomainSerialTarget{Port: &port},
	}}.ApplyTo(devs)

	if len(devs.Serials) != 1 {
		t.Fatalf("expected 1 serial, got %d", len(devs.Serials))
	}
	if devs.Serials[0].Target.Port == nil || *devs.Serials[0].Target.Port != 0 {
		t.Error("serial port should be 0")
	}
}

func TestConsole_ApplyTo(t *testing.T) {
	devs := emptyDeviceList()
	port := uint(0)
	Console{libvirtxml.DomainConsole{
		Source: &libvirtxml.DomainChardevSource{Pty: &libvirtxml.DomainChardevSourcePty{}},
		Target: &libvirtxml.DomainConsoleTarget{Type: "serial", Port: &port},
	}}.ApplyTo(devs)

	if len(devs.Consoles) != 1 {
		t.Fatalf("expected 1 console, got %d", len(devs.Consoles))
	}
	if devs.Consoles[0].Target.Type != "serial" {
		t.Errorf("console target type: expected serial, got %s", devs.Consoles[0].Target.Type)
	}
}

func TestGraphic_ApplyTo(t *testing.T) {
	devs := emptyDeviceList()
	Graphic{libvirtxml.DomainGraphic{
		VNC: &libvirtxml.DomainGraphicVNC{Port: -1, AutoPort: "yes"},
	}}.ApplyTo(devs)

	if len(devs.Graphics) != 1 {
		t.Fatalf("expected 1 graphic, got %d", len(devs.Graphics))
	}
	if devs.Graphics[0].VNC == nil {
		t.Fatal("vnc should not be nil")
	}
	if devs.Graphics[0].VNC.AutoPort != "yes" {
		t.Errorf("autoport: expected yes, got %s", devs.Graphics[0].VNC.AutoPort)
	}
}

func TestDevicerInterface_AllTypesImplement(t *testing.T) {
	devs := emptyDeviceList()
	devices := []Devicer{
		Disk{},
		NetworkInterface{},
		Channel{},
		Serial{},
		Console{},
		Graphic{},
	}
	for _, d := range devices {
		d.ApplyTo(devs)
	}
}
