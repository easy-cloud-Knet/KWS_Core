package conn

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"libvirt.org/go/libvirt"
)

func (SI *SystemInfo) GetInfo(domain *Domain) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
	}
	SI.Memory.Total = v.Total / 1024 / 1024 / 1024
	SI.Memory.Used = v.Used / 1024 / 1024 / 1024
	SI.Memory.Available = v.Available / 1024 / 1024 / 1024
	SI.Memory.UsedPercent = v.UsedPercent

	usage, err := disk.Usage("/")
	if err != nil {
		fmt.Println(err)
		return err
	}

	SI.Disks.Total = usage.Total / 1024 / 1024 / 1024
	SI.Disks.Used = usage.Used / 1024 / 1024 / 1024
	SI.Disks.Free = usage.Free / 1024 / 1024 / 1024
	SI.Disks.UsedPercent = usage.UsedPercent

	return nil
}

func (DI *DomainInfo) GetInfo(domain *Domain) error {
	info, err := domain.Domain.GetInfo()
	if err != nil {
		log.Println(err)
	}
	DI.State = info.State
	DI.MaxMem = info.MaxMem
	DI.Memory = info.Memory
	DI.NrVirtCpu = info.NrVirtCpu
	DI.CpuTime = info.CpuTime
	//basic info can be added
	return nil
}

func (DP *DomainState) GetInfo(domain *Domain) error {
	info, _, err := domain.Domain.GetState()
	//searching for coresponding second parameter, "Reason"
	if err != nil {
		log.Println(err)
		return err
	}

	uuidBytes, _ := domain.Domain.GetUUID()
	uuidParsed, _ := uuid.FromBytes(uuidBytes)

	DP.DomainState = info
	DP.UUID = string(uuidParsed.String())
	userInfo, err := domain.Domain.GetGuestInfo(libvirt.DOMAIN_GUEST_INFO_USERS, 0)
	if err != nil {
		log.Println(err)
		return err
	}
	DP.Users = userInfo.Users
	return nil
}

func DomainDetailFactory(Handler DataTypeHandler, Seeker DomainSeeker) *DomainDetail {
	return &DomainDetail{
		DataHandle:   make([]DataTypeHandler, 0),
		DomainSeeker: Seeker,
	}
}

func DataTypeRouter(types DomainDataType) (DataTypeHandler, error) {
	switch types {
	case DomState:
		return &DomainState{}, nil
	case BasicInfo:
		return &DomainInfo{}, nil
	case GuestInfoUser:
		return &DomainInfo{}, nil
	case GuestInfoOS:
		return &DomainInfo{}, nil
	case GuestInfoFS:
		return &DomainInfo{}, nil
	case GuestInfoDisk:
		return &DomainInfo{}, nil
	case HostInfo:
		return &SystemInfo{}, nil
	}
	return &DomainInfo{}, fmt.Errorf("not valid parameters for DomainDataType provided")
}



func (DSU *DomainSeekingByUUID) ReturnDomain() ([]*Domain, error) {
	fmt.Println((DSU))
	if len(DSU.Domain) == 0 {
		err :=DSU.SetDomain()
		if err!=nil{
			return nil, fmt.Errorf("no such Domain exists")
		}
	}
	return DSU.Domain, nil
}



func (DSS *DomainSeekingByStatus) ReturnDomain() ([]*Domain, error) {
	if len(DSS.DomList) == 0 {
		err :=DSS.SetDomain()
		if err!=nil{
			return nil, fmt.Errorf("no such Domain exists")
		}
	}
	return DSS.DomList, nil
}

func ReturnUUID(UUID string) (uuid.UUID, error) {
	uuidParsed, err := uuid.Parse(UUID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidParsed, nil
}

func (DSU *DomainSeekingByUUID) SetDomain() error {
	parsedUUID, err := uuid.Parse(DSU.UUID)
	if err != nil {
		return fmt.Errorf("invalid uuid format: %w", err)
	}
	domain, err := DSU.LibvirtInst.LookupDomainByUUID(parsedUUID[:])
	if err != nil {
		return fmt.Errorf("invalid uuid format: %w", err)
	}else if domain==nil {
		return fmt.Errorf("no such Domain For operation exists")
	}

	Dom := make([]*Domain, 1)
	Dom[0] = &Domain{
		Domain:      domain,
		DomainMutex: sync.Mutex{},
	}
	DSU.Domain = Dom
	return nil
}

func (DSS *DomainSeekingByStatus) SetDomain() error {
	doms, err := DSS.LibvirtInst.ListAllDomains(DSS.Status)
	if err != nil {
		fmt.Println("error while retrieving domain List with status")
		return fmt.Errorf("invalid uuid format: %w", err)
	}else if len(doms)==0{
		return fmt.Errorf("no such Domain For operation exists")
	}
	Domains := make([]*Domain, 0, len(doms))

	for i := range doms {
		Domains = append(Domains, &Domain{Domain: &doms[i], DomainMutex: sync.Mutex{}})
	}

	DSS.DomList = Domains
	return nil
}

func DomSeekStatusFactory(LibInstance *libvirt.Connect,flag libvirt.ConnectListAllDomainsFlags)*DomainSeekingByStatus{
	return &DomainSeekingByStatus{
		LibvirtInst: LibInstance,
		Status: flag,
		DomList: make([]*Domain, 0),
	}
}

func DomSeekUUIDFactory(LibInstance *libvirt.Connect,UUID string)*DomainSeekingByUUID{
	return &DomainSeekingByUUID{ 
		LibvirtInst: LibInstance,
		UUID:        UUID,
		Domain:      make([]*Domain, 0),
	}
}