package parsor

import (
	"fmt"
)
func(XP *VM_CREATE_XML) XML_Parsor(spec *VM_Init_Info) error {
	XP= &VM_CREATE_XML{
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
			Emulator: "/usr/bin/kvm",
			Disks: []Disk{
				{
					Type:   "file",
					Device: "disk",
					Driver: Driver{
						Name: "qemu",
						Type: "qcow2",
					},
					Source: Source{
						File: fmt.Sprintf("/var/lib/kws/%s/%s.qcow2", spec.UUID, spec.UUID),
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
						File: fmt.Sprintf("/var/lib/kws/%s/cidata.iso", spec.UUID),
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
			Interfaces: []Interface{
				{
					Type: "bridge",
					Source: NetworkSource{
						Bridge: "virbr1",
					},
					Model: InterfaceModel{
						Type: "virtio",
					},
				},
			},
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
 