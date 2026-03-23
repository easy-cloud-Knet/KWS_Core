package parsor

import (
	"encoding/xml"
	"strings"
	"testing"

	network "github.com/easy-cloud-Knet/KWS_Core/internal/net"
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
)

var testSpec = &vmtypes.VM_Init_Info{
	DomName: "test-vm",
	UUID:    "123e4567-e89b-12d3-a456-426614174000",
	OS:      "debian",
	HardwardInfo: vmtypes.HardwareInfo{
		CPU:    2,
		Memory: 4,
		Disk:   20,
	},
	NetConf: network.NetDefine{
		Ips: []string{"192.168.1.100"},
	},
	MacAddr: "52:54:00:ab:cd:ef",
	SDNUUID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
}

func TestBuildInterface_OVN(t *testing.T) {
	NetworkMode = "ovn"

	iface := buildInterface(testSpec)

	if iface.Source.Bridge != "br-int" {
		t.Errorf("bridge: expected br-int, got %s", iface.Source.Bridge)
	}
	if iface.Virtualport == nil {
		t.Fatal("virtualport is nil in ovn mode")
	}
	if iface.Virtualport.Type != "openvswitch" {
		t.Errorf("virtualport type: expected openvswitch, got %s", iface.Virtualport.Type)
	}
	if iface.Virtualport.Parameter.InterfaceID != testSpec.SDNUUID {
		t.Errorf("interfaceid: expected %s, got %s", testSpec.SDNUUID, iface.Virtualport.Parameter.InterfaceID)
	}
	if iface.MTU == nil {
		t.Fatal("mtu is nil in ovn mode")
	}
	if iface.MTU.Size != 1450 {
		t.Errorf("mtu size: expected 1450, got %d", iface.MTU.Size)
	}
}

func TestBuildInterface_Bridge(t *testing.T) {
	NetworkMode = "bridge"
	defer func() { NetworkMode = "ovn" }()

	iface := buildInterface(testSpec)

	if iface.Source.Bridge != "virbr0" {
		t.Errorf("bridge: expected virbr0, got %s", iface.Source.Bridge)
	}
	if iface.Virtualport != nil {
		t.Error("virtualport should be nil in bridge mode")
	}
	if iface.MTU != nil {
		t.Error("mtu should be nil in bridge mode")
	}
}

func TestBuildInterface_CommonFields(t *testing.T) {
	for _, mode := range []string{"ovn", "bridge"} {
		NetworkMode = mode
		iface := buildInterface(testSpec)

		if iface.Type != "bridge" {
			t.Errorf("[%s] type: expected bridge, got %s", mode, iface.Type)
		}
		if iface.MacAddress.Address != testSpec.MacAddr {
			t.Errorf("[%s] mac: expected %s, got %s", mode, testSpec.MacAddr, iface.MacAddress.Address)
		}
		if iface.Model.Type != "virtio" {
			t.Errorf("[%s] model: expected virtio, got %s", mode, iface.Model.Type)
		}
	}
	NetworkMode = "ovn"
}

func TestXMLParsor(t *testing.T) {
	NetworkMode = "ovn"

	var xp VM_CREATE_XML
	if err := xp.XML_Parsor(testSpec); err != nil {
		t.Fatal(err)
	}

	if xp.Name != testSpec.DomName {
		t.Errorf("name: expected %s, got %s", testSpec.DomName, xp.Name)
	}
	if xp.UUID != testSpec.UUID {
		t.Errorf("uuid: expected %s, got %s", testSpec.UUID, xp.UUID)
	}
	if xp.Memory.Size != testSpec.HardwardInfo.Memory {
		t.Errorf("memory: expected %d, got %d", testSpec.HardwardInfo.Memory, xp.Memory.Size)
	}
	if xp.VCPU.Count != testSpec.HardwardInfo.CPU {
		t.Errorf("vcpu: expected %d, got %d", testSpec.HardwardInfo.CPU, xp.VCPU.Count)
	}
	if len(xp.Devices.Disks) != 2 {
		t.Errorf("disks: expected 2, got %d", len(xp.Devices.Disks))
	}
	if len(xp.Devices.Interfaces) != 1 {
		t.Errorf("interfaces: expected 1, got %d", len(xp.Devices.Interfaces))
	}
}

func TestXMLParsor_DiskPaths(t *testing.T) {
	var xp VM_CREATE_XML
	if err := xp.XML_Parsor(testSpec); err != nil {
		t.Fatal(err)
	}

	uuid := testSpec.UUID
	wantQcow := "/var/lib/kws/" + uuid + "/" + uuid + ".qcow2"
	wantISO := "/var/lib/kws/" + uuid + "/cidata.iso"

	if xp.Devices.Disks[0].Source.File != wantQcow {
		t.Errorf("disk path: expected %s, got %s", wantQcow, xp.Devices.Disks[0].Source.File)
	}
	if xp.Devices.Disks[1].Source.File != wantISO {
		t.Errorf("iso path: expected %s, got %s", wantISO, xp.Devices.Disks[1].Source.File)
	}
	if xp.Devices.Disks[1].ReadOnly == nil {
		t.Error("cdrom should be readonly")
	}
}

func TestXMLParsor_MarshalOVN(t *testing.T) {
	NetworkMode = "ovn"

	var xp VM_CREATE_XML
	if err := xp.XML_Parsor(testSpec); err != nil {
		t.Fatal(err)
	}

	out, err := xml.MarshalIndent(xp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	xmlStr := string(out)

	for _, want := range []string{
		`<domain type="kvm">`,
		testSpec.DomName,
		testSpec.UUID,
		`<virtualport type="openvswitch">`,
		testSpec.SDNUUID,
		`<mtu size="1450"`,
		`br-int`,
	} {
		if !strings.Contains(xmlStr, want) {
			t.Errorf("XML missing %q", want)
		}
	}
}

func TestXMLParsor_MarshalBridge(t *testing.T) {
	NetworkMode = "bridge"
	defer func() { NetworkMode = "ovn" }()

	var xp VM_CREATE_XML
	if err := xp.XML_Parsor(testSpec); err != nil {
		t.Fatal(err)
	}

	out, err := xml.MarshalIndent(xp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	xmlStr := string(out)

	if !strings.Contains(xmlStr, "virbr0") {
		t.Error("XML missing virbr0")
	}
	if strings.Contains(xmlStr, "virtualport") {
		t.Error("XML should not contain virtualport in bridge mode")
	}
	if strings.Contains(xmlStr, "<mtu") {
		t.Error("XML should not contain mtu in bridge mode")
	}
}
