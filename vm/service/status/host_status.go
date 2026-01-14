package status

import (
	"errors"
	"fmt"
	"time"

	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domain_status"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func (CI *HostCpuInfo) GetHostInfo(status *domStatus.DomainListStatus) error {
	t, err := cpu.Times(false) //time
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, err)
	}

	p, err := cpu.Percent(time.Second, false)
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, err)
	}

	if len(t) > 0 {
		CI.System = t[0].System
		CI.Idle = t[0].Idle
	}
	CI.Usage = p[0]

	CI.Desc.EmitStatus(status)



	return nil
}

func (MI *HostMemoryInfo) GetHostInfo(status *domStatus.DomainListStatus) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, err)

	}

	MI.Total = v.Total / 1024 / 1024 / 1024
	MI.Used = v.Used / 1024 / 1024 / 1024
	MI.Available = v.Available / 1024 / 1024 / 1024
	MI.UsedPercent = v.UsedPercent

	return nil
}

func (HDI *HostDiskInfo) GetHostInfo(status *domStatus.DomainListStatus) error {
	d, err := disk.Usage("/")
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, err)

	}

	HDI.Total = d.Total / 1024 / 1024 / 1024
	HDI.Used = d.Used / 1024 / 1024 / 1024
	HDI.Free = d.Free / 1024 / 1024 / 1024
	HDI.UsedPercent = d.UsedPercent

	return nil
}

func (SI *HostGeneralInfo) GetHostInfo(status *domStatus.DomainListStatus) error {
	err:=SI.CPU.GetHostInfo(status)
	if err!=nil{
		return virerr.ErrorGen(virerr.HostStatusError, fmt.Errorf("general Status:error retreving host Status %w",err))
	}
	err=SI.Disk.GetHostInfo(status)
	if err!=nil{
		return virerr.ErrorGen(virerr.HostStatusError, fmt.Errorf("general Status:error retreving host Status %w",err))
	}
	err=SI.Memory.GetHostInfo(status)
	if err!=nil{
		return virerr.ErrorGen(virerr.HostStatusError, fmt.Errorf("general Status:error retreving host Status %w",err))
	}

	return nil
}





func (HSI *HostSystemInfo) GetHostInfo(status *domStatus.DomainListStatus) error {
	u, err := host.Uptime()
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, err)

	}

	b, err := host.BootTime()
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, err)

	}

	HSI.Uptime = u
	HSI.BootTime = b

	temp, err := host.SensorsTemperatures()
	if err == nil {
		for _, t := range temp {
			if t.SensorKey == "coretemp" || t.SensorKey == "cpu" {
				HSI.CPU_Temp = t.Temperature
			} else if t.SensorKey == "dimm" || t.SensorKey == "ram" {
				HSI.RAM_Temp = t.Temperature
			}
		}
	}// 에러가 발생했을때 대처가 부족한거 같음

	return nil
}

func HostDataTypeRouter(types HostDataType) (HostDataTypeHandler, error) {
	// implemeantation of factory pattern
	// can be extened as DI pattern if there are more complex dependencies

	switch types {
	case CpuInfo:
		return &HostCpuInfo{
			Desc: &domStatus.VCPUStatus{},
		}, nil
	case MemInfo:
		return &HostMemoryInfo{
			Desc: &domStatus.VCPUStatus{},// --- IGNORE ---
		}, nil
	case DiskInfoHi:
		return &HostDiskInfo{
			Desc: &domStatus.VCPUStatus{},// --- IGNORE ---
		}, nil
	case SystemInfoHi:
		return &HostSystemInfo{}, nil
	case GeneralInfo:
		return &HostGeneralInfo{},nil
	}
		
		return nil, virerr.ErrorGen(virerr.HostStatusError, errors.New("not valid parameters for HostDataType provided"))
}



func HostInfoHandler(handler HostDataTypeHandler, status *domStatus.DomainListStatus) (*HostDetail, error) {
	if err := handler.GetHostInfo(status); err != nil {
		return nil, err
	}
	return &HostDetail{
		HostDataHandle: handler,
	}, nil
}

