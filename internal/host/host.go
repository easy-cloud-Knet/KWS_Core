package host

import (
	"errors"
	"fmt"
	"time"

	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domainList_status"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type DataTypeHandler interface {
	GetHostInfo(*domStatus.DomainListStatus) error
}

type DataType uint

const (
	CPU DataType = iota
	Memory
	Disk
	System
	General
	DomainAll
)

type Detail struct {
	DataHandle DataTypeHandler
}

type CpuInfo struct {
	System float64               `json:"system_time"`
	Idle   float64               `json:"idle_time"`
	Usage  float64               `json:"usage_percent"`
	Desc   *domStatus.VCPUStatus `json:"vcpu_status"`
}

type MemoryInfo struct {
	Total          uint64                `json:"total_gb"`
	Used           uint64                `json:"used_gb"`
	Available      uint64                `json:"available_gb"`
	UsedPercent    float64               `json:"used_percent"`
	ReservedMemory uint64                `json:"reservedmem"`
	Desc           *domStatus.VCPUStatus `json:"vcpu_status"`
}

type DiskInfo struct {
	Total       uint64                `json:"total_gb"`
	Used        uint64                `json:"used_gb"`
	Free        uint64                `json:"free_gb"`
	UsedPercent float64               `json:"used_percent"`
	Desc        *domStatus.VCPUStatus `json:"vcpu_status"`
}

type SystemInfo struct {
	Uptime   uint64  `json:"uptime_seconds"`
	BootTime uint64  `json:"boot_time_epoch"`
	CPU_Temp float64 `json:"cpu_temperature,omitempty"`
	RAM_Temp float64 `json:"ram_temperature,omitempty"`
}

type GeneralInfo struct {
	CPU    CpuInfo    `json:"cpuInfo"`
	Memory MemoryInfo `json:"memoryInfo"`
	Disk   DiskInfo   `json:"DiskInfo"`
}

func (CI *CpuInfo) GetHostInfo(status *domStatus.DomainListStatus) error {
	t, err := cpu.Times(false)
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

func (MI *MemoryInfo) GetHostInfo(_ *domStatus.DomainListStatus) error {
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

func (HDI *DiskInfo) GetHostInfo(_ *domStatus.DomainListStatus) error {
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

func (SI *GeneralInfo) GetHostInfo(status *domStatus.DomainListStatus) error {
	if err := SI.CPU.GetHostInfo(status); err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, fmt.Errorf("general Status:error retreving host Status %w", err))
	}
	if err := SI.Disk.GetHostInfo(status); err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, fmt.Errorf("general Status:error retreving host Status %w", err))
	}
	if err := SI.Memory.GetHostInfo(status); err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, fmt.Errorf("general Status:error retreving host Status %w", err))
	}
	return nil
}

func (HSI *SystemInfo) GetHostInfo(_ *domStatus.DomainListStatus) error {
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
	}
	return nil
}

func DataTypeRouter(t DataType) (DataTypeHandler, error) {
	switch t {
	case CPU:
		return &CpuInfo{Desc: &domStatus.VCPUStatus{}}, nil
	case Memory:
		return &MemoryInfo{Desc: &domStatus.VCPUStatus{}}, nil
	case Disk:
		return &DiskInfo{Desc: &domStatus.VCPUStatus{}}, nil
	case System:
		return &SystemInfo{}, nil
	case General:
		return &GeneralInfo{}, nil
	}
	return nil, virerr.ErrorGen(virerr.HostStatusError, errors.New("not valid parameters for HostDataType provided"))
}

func InfoHandler(handler DataTypeHandler, status *domStatus.DomainListStatus) (*Detail, error) {
	if err := handler.GetHostInfo(status); err != nil {
		return nil, err
	}
	return &Detail{DataHandle: handler}, nil
}
