package process

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/process"
)

type ProcessInfo struct {
	PID        int32
	Name       string
	CpuPercent float64
}

func getProcessInfo() []ProcessInfo {

	processes, err := process.Processes()
	if err != nil {
		fmt.Println(" Error fetching processes: ", err)
		return nil
	}

	var processInfos []ProcessInfo

	for _, p := range processes {
		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue
		}

		name, err := p.Name()
		if err != nil {
			name = "UnKnown"
		}

		processInfos = append(processInfos, ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			CpuPercent: cpuPercent,
		})
	}
	return processInfos
}

func Update() []ProcessInfo {
	pInfo := getProcessInfo()
	top := getTopProcess(pInfo)
	return top
}

func getTopProcess(pInfo []ProcessInfo) []ProcessInfo {
	sort.Slice(pInfo, func(i, j int) bool {
		return pInfo[i].CpuPercent > pInfo[j].CpuPercent
	})

	for i, p := range pInfo {
		if i >= 5 {
			break
		}
		fmt.Printf("%-8d %-20s %.2f%%\n", p.PID, p.Name, p.CpuPercent)
	}

	return pInfo
}
