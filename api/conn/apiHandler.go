package conn

import (
	"fmt"
	"net/http"
	"os"

	"libvirt.org/go/libvirt"
)





func (i * InstHandler) ReturnStatus(w http.ResponseWriter,r * http.Request){
	fmt.Println("getStatus request income")

	i.ReturnDomainNameList(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
}

func (i *InstHandler) CreateDomainWithXML(w http.ResponseWriter, r *http.Request) {

	// 파일 포인터를 슬라이스에 담습니다.
	files := []os.File{}
	xmlConfig := `
		<domain type='kvm'>
			<name>demo2</name>
			<uuid>4dea24b3-1d52-d8f3-2516-782e98a23fa0</uuid>
			<memory>131072</memory>
			<vcpu>1</vcpu>
			<os>
				<type arch="x86_64">hvm</type>
			</os>
			<clock sync="localtime"/>
			<devices>
				<emulator>/usr/bin/qemu-kvm</emulator>
				<disk type='file' device='disk'>
					<source file='/var/lib/libvirt/images/demo2.img'/>
					<target dev='hda'/>
				</disk>
				<interface type='network'>
					<source network='default'/>
					<mac address='24:42:53:21:52:45'/>
				</interface>
				<graphics type='vnc' port='-1' keymap='de'/>
			</devices>
		</domain>
	`

	// 추가 파일이 없는 경우 빈 슬라이스를 전달합니다.

	// DomainCreateXMLWithFiles를 호출하여 도메인을 생성합니다.
	domain, err := i.LibvirtInst.DomainCreateXMLWithFiles(xmlConfig, files, libvirt.DOMAIN_START_FORCE_BOOT)
	if err != nil {
		http.Error(w, "Failed to create domain", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Domain created: %v", domain)
}
