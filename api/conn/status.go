package conn

import (
	"errors"
	"log"
	"sync"

	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"libvirt.org/go/libvirt"
)

func (SI *SystemInfo) GetInfo(domain *Domain) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
		return virerr.ErrorGen(virerr.HostStatusError,err)
	}
	SI.Memory.Total = v.Total / 1024 / 1024 / 1024
	SI.Memory.Used = v.Used / 1024 / 1024 / 1024
	SI.Memory.Available = v.Available / 1024 / 1024 / 1024
	SI.Memory.UsedPercent = v.UsedPercent

	usage, err := disk.Usage("/")
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, err)
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
		return virerr.ErrorGen(virerr.DomainStatusError, err)
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
		return virerr.ErrorGen(virerr.DomainStatusError, err)
	}

	uuidBytes,err := domain.Domain.GetUUID()
	if err!= nil{
		return virerr.ErrorGen(virerr.InvalidUUID, err)
	}
	uuidParsed, err := uuid.FromBytes(uuidBytes)
	if err!= nil{
		return virerr.ErrorGen(virerr.InvalidUUID, err)
	}

	DP.DomainState = info
	DP.UUID = string(uuidParsed.String())
	userInfo, err := domain.Domain.GetGuestInfo(libvirt.DOMAIN_GUEST_INFO_USERS, 0)
	if err != nil {
		log.Println("error retreving guest info")
		return err
	}
	DP.Users = userInfo.Users
	return nil
}

func DomainDetailFactory(Handler DataTypeHandler, dom *Domain) *DomainDetail {
	return &DomainDetail{
		DataHandle: Handler,
		Domain: dom,
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
	return nil, virerr.ErrorGen(virerr.InvalidParameter, errors.New("invalid flag for DataRoute entereed "))
}



func (DSU *DomainSeekingByUUID) ReturnDomain() (*Domain, error) {
	parsedUUID, err := uuid.Parse(DSU.UUID)
	if err != nil {
		return nil,virerr.ErrorGen(virerr.InvalidUUID, err)
	}
	domain, err := DSU.LibvirtInst.LookupDomainByUUID(parsedUUID[:])
	if err != nil {
		return nil,virerr.ErrorGen(virerr.DomainSearchError, err)
	}else if domain==nil {
		return nil,virerr.ErrorGen(virerr.NoSuchDomain, err)
	}

	
	return &Domain{
		Domain:      domain,
		domainMutex: sync.Mutex{},
	}, nil
}



 
func ReturnUUID(UUID string) (*uuid.UUID, error) {
	uuidParsed, err := uuid.Parse(UUID)
	if err != nil {
		return nil, err
	}
	return &uuidParsed, nil
}



func DomSeekUUIDFactory(LibInstance *libvirt.Connect,UUID string)*DomainSeekingByUUID{
	return &DomainSeekingByUUID{ 
		LibvirtInst: LibInstance,
		UUID:        UUID,
	}
}