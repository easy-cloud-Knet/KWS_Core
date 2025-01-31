package service

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func GetSystemInfo() (string, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "", fmt.Errorf("memory error: %v", err)
	}

	memoryInfo := MemoryInfo{
		Total:       v.Total / 1024 / 1024 / 1024,
		Used:        v.Used / 1024 / 1024 / 1024,
		Available:   v.Available / 1024 / 1024 / 1024,
		UsedPercent: v.UsedPercent,
	}

	partitions, err := disk.Partitions(false)
	if err != nil {
		return "", fmt.Errorf("disk error: %v", err)
	}

	var disks []DiskInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		disks = append(disks, DiskInfo{
			Mountpoint:  p.Mountpoint,
			Total:       usage.Total / 1024 / 1024 / 1024,
			Used:        usage.Used / 1024 / 1024 / 1024,
			Free:        usage.Free / 1024 / 1024 / 1024,
			UsedPercent: usage.UsedPercent,
		})
	}

	systemInfo := SystemInfo{
		Memory: memoryInfo,
		Disks:  disks,
	}

	jsonData, err := json.MarshalIndent(systemInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}

	return string(jsonData), nil
}
