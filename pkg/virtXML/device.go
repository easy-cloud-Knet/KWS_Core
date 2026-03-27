package virtxml

import libvirtxml "libvirt.org/libvirt-go-xml"

// Devicer는 libvirtxml.DomainDeviceList에 자신을 추가할 수 있는 디바이스.
type Devicer interface {
	ApplyTo(devs *libvirtxml.DomainDeviceList)
}

type Disk struct{ libvirtxml.DomainDisk }

func (d Disk) ApplyTo(devs *libvirtxml.DomainDeviceList) {
	devs.Disks = append(devs.Disks, d.DomainDisk)
}

type NetworkInterface struct{ libvirtxml.DomainInterface }

func (i NetworkInterface) ApplyTo(devs *libvirtxml.DomainDeviceList) {
	devs.Interfaces = append(devs.Interfaces, i.DomainInterface)
}

type Channel struct{ libvirtxml.DomainChannel }

func (c Channel) ApplyTo(devs *libvirtxml.DomainDeviceList) {
	devs.Channels = append(devs.Channels, c.DomainChannel)
}

type Serial struct{ libvirtxml.DomainSerial }

func (s Serial) ApplyTo(devs *libvirtxml.DomainDeviceList) {
	devs.Serials = append(devs.Serials, s.DomainSerial)
}

type Console struct{ libvirtxml.DomainConsole }

func (c Console) ApplyTo(devs *libvirtxml.DomainDeviceList) {
	devs.Consoles = append(devs.Consoles, c.DomainConsole)
}

type Graphic struct{ libvirtxml.DomainGraphic }

func (g Graphic) ApplyTo(devs *libvirtxml.DomainDeviceList) {
	devs.Graphics = append(devs.Graphics, g.DomainGraphic)
}
