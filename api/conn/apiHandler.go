package conn

import (
	"fmt"
	"net/http"

	"log"

	"encoding/json"

	"libvirt.org/go/libvirt"
)





func (i * InstHandler) ReturnStatus(w http.ResponseWriter,r * http.Request){
	fmt.Println("getStatus request income")

	Domlist,_:= i.ReturnDomainNameList(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	encoder := json.NewEncoder(w)
	encoder.Encode(&Domlist)

}




func (i *InstHandler) CreateDomainWithXML(w http.ResponseWriter, r *http.Request) {
	// 파일 포인터를 슬라이스에 담습니다.
	xmlConfig := `<domain type='kvm'>
  <name>cloud-vm</name>
  <memory unit='GiB'>2</memory>
  <vcpu placement='static'>2</vcpu>
  <os>
    <type arch='x86_64' >hvm</type>
    <boot dev='hd'/>
  </os>
  <devices>
    <features>
	  <acpi/>
	</features>
    <emulator>/usr/bin/kvm</emulator>
    <disk type='file' device='disk'>
      <driver name='qemu' type='qcow2'/>
      <source file='/var/lib/libvirt/vmList/user2/ubuntuInst.qcow2'/>
      <target dev='vda' bus='virtio'/>
    </disk>
    <disk type='file' device='cdrom'>
      <driver name='qemu' type='raw'/>
      <source file='/var/lib/libvirt/vmList/user2/cidata.iso'/>
      <target dev='hda' bus='ide'/>
      <readonly/>
    </disk>
   <serial type='pty'>
		  <target port='0'/>
		</serial>
		<console type='pty'>
		  <target type='serial' port='0'/>
		</console>
	<interface type='network'>
      <source network='default'/>
      <model type='virtio'/>
    </interface>
    <graphics type='vnc' port='-1' autoport='yes'/>
  </devices>
</domain>
	`
	

	// 추가 파일이 없는 경우 빈 슬라이스를 전달합니다.

	// DomainCreateXMLWithFiles를 호출하여 도메인을 생성합니다.
	domain, err := i.LibvirtInst.DomainCreateXML(xmlConfig, libvirt.DOMAIN_NONE)
	if err != nil {
		log.Fatal(err)
		
	}

	fmt.Fprintf(w, "Domain created: %v", domain)
}
