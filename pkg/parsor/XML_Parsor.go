package parsor

import (
	"fmt"

	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
)

// TODO: add MarshalXML() ([]byte, error) method to VM_CREATE_XML
// to enable XMLDefiner interface to cover marshaling, allowing full mock in creation tests.

// NetworkMode is injected at build time via -ldflags. Default: "ovn".
// Supported values: "ovn", "bridge"
var NetworkMode = "ovn"

func buildInterface(spec *vmtypes.VM_Init_Info) Interface {
	iface := Interface{
		Type:       "bridge",
		MacAddress: MacAddress{Address: spec.MacAddr},
		Model:      InterfaceModel{Type: "virtio"},
	}
	if NetworkMode == "ovn" {
		iface.Source = NetworkSource{Bridge: "br-int"}
		iface.Virtualport = &VirPort{
			Type:      "openvswitch",
			Parameter: Parameter{InterfaceID: spec.SDNUUID},
		}
		iface.MTU = &MTU{Size: 1450}
	} else {
		iface.Source = NetworkSource{Bridge: "virbr0"}
	}
	return iface
}

func (XP *VM_CREATE_XML) XML_Parsor(spec *vmtypes.VM_Init_Info) error {
	*XP = VM_CREATE_XML{
		Type: "kvm",
		Name: spec.DomName,
		UUID: spec.UUID,
		Memory: Memory{
			Unit: "GiB",
			Size: spec.HardwardInfo.Memory,
		},
		VCPU: VCPU{
			Placement: "static",
			Count:     spec.HardwardInfo.CPU,
		},
		Features: Features{
			ACPI: ACPI{},
		},
		OS: OS{
			Type: OSType{
				Arch: "x86_64",
				Type: "hvm",
			},
			Boot: Boot{
				Dev: "hd",
			},
		},
		Devices: Devices{
			Emulator: config.EmulatorPath,
			Disks: []Disk{
				{
					Type:   "file",
					Device: "disk",
					Driver: Driver{
						Name: "qemu",
						Type: "qcow2",
					},
					Source: Source{
						File: fmt.Sprintf("%s/%s/%s.qcow2", config.StorageBase, spec.UUID, spec.UUID),
					},
					Target: Target{
						Dev: "vda",
						Bus: "virtio",
					},
				},
				{
					Type:   "file",
					Device: "cdrom",
					Driver: Driver{
						Name: "qemu",
						Type: "raw",
					},
					Source: Source{
						File: fmt.Sprintf("%s/%s/cidata.iso", config.StorageBase, spec.UUID),
					},
					Target: Target{
						Dev: "sda",
						Bus: "ide",
					},
					ReadOnly: &ReadOnly{},
				},
			},
			Serial: Serial{
				Type: "pty",
				Target: TargetPort{
					Port: 0,
				},
			},
			Console: Console{
				Type: "pty",
				Target: ConsoleTarget{
					Type: "serial",
					Port: 0,
				},
			},
			Interfaces: []Interface{buildInterface(spec)},
			Graphics: Graphics{
				Type:     "vnc",
				Port:     -1,
				AutoPort: "yes",
			},
			Channels: []Channel{ // QEMU Guest Agent 설정 추가
				{
					Type: "unix",
					Source: ChannelSource{
						Mode: "bind",
					},
					Target: ChannelTarget{
						Type: "virtio",
						Name: "org.qemu.guest_agent.0",
					},
				},
			},
		},
	}

	return nil
}
