package conn

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func (CI *HostCpuInfo) GetHostInfo() error {
	t, err := cpu.Times(false) //time
	if err != nil {
		log.Println(err)
		return ErrorGen(HostStatusError, err)
	}

	p, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println(err)
		return ErrorGen(HostStatusError, err)
	}

	if len(t) > 0 {
		CI.System = t[0].System
		CI.Idle = t[0].Idle
	}
	CI.Usage = p[0]

	return nil
}

func (MI *HostMemoryInfo) GetHostInfo() error {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
		return ErrorGen(HostStatusError, err)

	}

	MI.Total = v.Total / 1024 / 1024 / 1024
	MI.Used = v.Used / 1024 / 1024 / 1024
	MI.Available = v.Available / 1024 / 1024 / 1024
	MI.UsedPercent = v.UsedPercent

	return nil
}

func (HDI *HostDiskInfo) GetHostInfo() error {
	d, err := disk.Usage("/")
	if err != nil {
		log.Println(err)
		return ErrorGen(HostStatusError, err)

	}

	HDI.Total = d.Total / 1024 / 1024 / 1024
	HDI.Used = d.Used / 1024 / 1024 / 1024
	HDI.Free = d.Free / 1024 / 1024 / 1024
	HDI.UsedPercent = d.UsedPercent

	return nil
}

func (HSI *HostSystemInfo) GetHostInfo() error {
	u, err := host.Uptime()
	if err != nil {
		log.Println(err)
		return ErrorGen(HostStatusError, err)

	}

	b, err := host.BootTime()
	if err != nil {
		log.Println(err)
		return ErrorGen(HostStatusError, err)

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

	switch types {
	case CpuInfo:
		return &HostCpuInfo{}, nil
	case MemInfo:
		return &HostMemoryInfo{}, nil
	case DiskInfoHi:
		return &HostDiskInfo{}, nil
	case SystemInfoHi:
		return &HostSystemInfo{}, nil
	}
		
		return nil, ErrorGen(HostStatusError, errors.New("not valid parameters for HostDataType provided"))
}



func HostDetailFactory(handler HostDataTypeHandler) (*HostDetail, error) {
	if err := handler.GetHostInfo(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &HostDetail{
		HostDataHandle: handler,
	}, nil
}

