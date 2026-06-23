package snapshotlibvirt

import "libvirt.org/go/libvirt"

// RegisterExternalSnapshot registers an external snapshot using pre-existing overlay files
// created by qemu-img. REUSE_EXT tells libvirt to reuse the files rather than creating new ones.
func RegisterExternalSnapshot(domain *libvirt.Domain, snapshotXML string) (*libvirt.DomainSnapshot, error) {
	flags := libvirt.DOMAIN_SNAPSHOT_CREATE_DISK_ONLY | libvirt.DOMAIN_SNAPSHOT_CREATE_REUSE_EXT
	return domain.CreateSnapshotXML(snapshotXML, flags)
}

func UpdateDeviceConfig(domain *libvirt.Domain, deviceXML string) error {
	return domain.UpdateDeviceFlags(deviceXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG)
}
