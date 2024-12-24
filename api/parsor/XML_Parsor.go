package parsor

import (
	"encoding/xml"
	"fmt"
)

//
func XML_Parsor() {
	vm := &VM_CREATE_XML{
		Type: "kvm",
		Name: "cloud-vm",
		UUID: "6a21d302-e2b0-4a53-a9a5-4b08021cbba2",
		Memory: Memory{
			Unit: "GiB",
			Size: 2,
		},
		VCPU: VCPU{
			Placement: "static",
			Count:     2,
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
						File: "/var/lib/kws/user1/user1.qcow2",
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
						File: "/var/lib/kws/user1/cidata.iso",
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
					Type: "network",
					Source: NetworkSource{
						Network: "default",
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
		},
	}

	output, _ := xml.MarshalIndent(vm, "", "  ")
	fmt.Println(string(output))
}
