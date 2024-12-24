package conn

import (
	"fmt"
	"net/http"

	"log"

	"encoding/json"
	"github.com/easy-cloud-knet/KWS_Core.git/api/parsor"
	"libvirt.org/go/libvirt"
)





func (i * InstHandler) ReturnStatus(w http.ResponseWriter,r * http.Request){
	fmt.Println("getStatus request income")

	Domlist,_:= i.ReturnDomainNameList(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	encoder := json.NewEncoder(w)
	encoder.Encode(&Domlist)

}




func (i *InstHandler) CreateDomainWithXML(config *parsor.VM_CREATE_XML) (*libvirt.Domain, error) {
	// 파일 포인터를 슬라이스에 담습니다.

	/*`
	<domain type='kvm'>
    <name>cloud-vm</name>
    <uuid>6a21d302-e2b0-4a53-a9a5-4b08021cbba2</uuid>
    <memory unit='GiB'>2</memory>
    <vcpu placement='static'>2</vcpu>
    <features>
      <acpi/>
    </features>
    <os>
      <type arch='x86_64'>hvm</type>
      <boot dev='hd'/>
    </os>
  <devices>
    <emulator>/usr/bin/kvm</emulator>
    <disk type='file' device='disk'>
      <driver name='qemu' type='qcow2'/>
      <source file='/var/lib/kws/user1/user1.qcow2'/>
      <target dev='vda' bus='virtio'/>
    </disk>
    <disk type='file' device='cdrom'>
      <driver name='qemu' type='raw'/>
      <source file='/var/lib/kws/user1/cidata.iso'/>
      <target dev='sda' bus='ide'/>
      <readonly/>
    </disk>
    <serial type='pty'>
      <target port='0'/>
    </serial>
    <console type='pty'>
      <target type='serial' port='0'/>
    </console>
    <interface type='bridge'>
      <source bridge='virbr1'/>
      <model type='virtio'/>
    </interface>
    <graphics type='vnc' port='-1' autoport='yes'/>
  </devices>
</domain>	`
	*/

	// 추가 파일이 없는 경우 빈 슬라이스를 전달합니다.

	// DomainCreateXMLWithFiles를 호출하여 도메인을 생성합니다.
	domain, err := i.LibvirtInst.DomainDefineXML(config)
	if err != nil {
		log.Fatal(err)
		
	}


  return domain ,err

}
