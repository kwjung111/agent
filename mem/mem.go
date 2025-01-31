package mem

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryCollector struct{}

type MemoryInfo struct {
	Total     uint64
	Available uint64
	Used      uint64
}

func getMem() (MemoryInfo, error) {
	virtMem, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}, err
	}
	fmt.Println("")
	ret, err := parse(virtMem)
	if err != nil {
		return MemoryInfo{}, err
	}
	return ret, nil
}

func parse(virtMem *mem.VirtualMemoryStat) (MemoryInfo, error) {
	total := virtMem.Total
	available := virtMem.Available
	used := virtMem.Used

	memoryInfo := MemoryInfo{
		Total:     total,
		Available: available,
		Used:      used,
	}

	return memoryInfo, nil
}

func (m *MemoryCollector) Update() (interface{}, error) {
	mem, err := getMem()
	if err != nil {
		return MemoryInfo{}, nil
	}
	return mem, nil
}
