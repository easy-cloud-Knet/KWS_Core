package conn

import (
	"fmt"
	"sync"

	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"libvirt.org/go/libvirt"
)


func DomGen(Dom *libvirt.Domain) *Domain{
	return &Domain{
		domainMutex: sync.Mutex{},
		Domain: Dom,
	}
}




func DomListConGen() *DomListControl{
	return &DomListControl{
		domainMutex: sync.Mutex{},
		DomainList: make(map[string]*Domain),
	}
}



func (DC *DomListControl) AddNewDomain(domain *Domain, uuid string){
	DC.domainMutex.Lock()
	DC.DomainList[uuid]= domain
	defer DC.domainMutex.Unlock()
	// 아직은 단순한 정도의 mutex만 구현, domainList 와 Domain이
	// 얼마나 복잡하냐에 따라 수정 가능.
}


func (DC *DomListControl) GetDomain(uuid string, LibvirtInst *libvirt.Connect)(*Domain, error){
	domain, Exist := DC.DomainList[uuid]; 
	if !Exist{
		DomainSeeker:= DomSeekUUIDFactory(LibvirtInst, uuid)
		domList,err :=DomainSeeker.ReturnDomain()
		if err!=nil{
			return nil, virerr.ErrorGen(virerr.NoSuchDomain,fmt.Errorf("no such domain exists with uuid of %s , %w", uuid,err))
		}
		domain=domList
		DC.AddNewDomain(domain,uuid)
		return domain, nil
	}

	return domain,nil
}

func (DC *DomListControl) DeleteDomain(uuid string, LibvirtInst *libvirt.Connect)(error){
	domain, Exist := DC.DomainList[uuid]; 
	if !Exist{
		DomainSeeker:= DomSeekUUIDFactory(LibvirtInst, uuid)
		dom,err :=DomainSeeker.ReturnDomain()
		if err!=nil{
			return virerr.ErrorGen(virerr.NoSuchDomain,fmt.Errorf("domain trying to delete already empty, uuid of %s , %w", uuid,err))
		}
		dom.Domain.Free()
		fmt.Println(dom)
		//도메인 삭제 로직 추가, 로그 추가 <--- Map과 싱크가 안맞았다는 얘기(골치 아픔)
		return nil
	}
	domain.Domain.Free()
	delete(DC.DomainList, uuid)
	return nil
}


func (DC *DomListControl) RetreiveAllDomain(LibvirtInst *libvirt.Connect)(error){
	domActive, err := LibvirtInst.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err!=nil{	
		fmt.Print("do something")
	}
	for i:= range domActive{
		uuid, err := domActive[i].GetUUID()
		if err!=nil{
			fmt.Println(err)
		}
		uuidStr := fmt.Sprintf("%x", uuid)
		DC.DomainList[uuidStr]= &Domain{
			Domain: &domActive[i],
			domainMutex: sync.Mutex{},
		}
	}
	domInactive, err := LibvirtInst.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err!=nil{
		fmt.Print("do something")

	}

	for i:= range domInactive{
		uuid, err := domInactive[i].GetUUID()
		if err!=nil{
			fmt.Println(err)
		}
		uuidStr := fmt.Sprintf("%x", uuid)
		DC.DomainList[uuidStr]= &Domain{
			Domain: &domInactive[i],
			domainMutex: sync.Mutex{},
		}
	}
	fmt.Println(DC.DomainList)


	return nil
}