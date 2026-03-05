package config

var (
	StorageBase        = "/var/lib/kws"
	LogDir             = "/var/log/kws/"
	LibvirtURI         = "qemu:///system"
	ServerPort         = 8080
	DefaultDNS         = "8.8.8.8"
	NetworkBridge      = "br-int"
	NetworkMTU         = 1450
	NetworkVirtualPort = "openvswitch"
	EmulatorPath       = "/usr/bin/kvm"
)
