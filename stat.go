package main

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	Percent struct {
		CPU  float64 `json:"cpu"`
		Mem  float64 `json:"mem"`
		Swap float64 `json:"swap"`
		Disk float64 `json:"disk"`
	} `json:"percent"`
	Load   *load.AvgStat `json:"load"`
	Memory UsageStat     `json:"memory"`
	Swap   UsageStat     `json:"swap"`
	Disk   struct {
		UsageStat
		Read  uint64 `json:"read"`
		Write uint64 `json:"write"`
	} `json:"disk"`
	Network struct {
		Rx   uint64 `json:"rx"`
		Tx   uint64 `json:"tx"`
		In   uint64 `json:"in"`
		Out  uint64 `json:"out"`
		Min  uint64 `json:"min"`
		Mout uint64 `json:"mout"`
	} `json:"network"`
	Uptime   uint64 `json:"uptime"`
	Hostname string `json:"hostname"`
}

type UsageStat struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
}

func getSystemInfo() (result SystemInfo, err error) {
	result.Percent.CPU = backgroundStat.CPUPercent

	result.Load, err = load.Avg()
	if err != nil {
		return result, err
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return result, err
	}
	result.Memory = UsageStat{
		Total: vm.Total,
		Used:  vm.Used,
		Free:  vm.Free,
	}
	result.Percent.Mem = vm.UsedPercent

	sm, err := mem.SwapMemory()
	if err != nil {
		return result, err
	}
	result.Swap = UsageStat{
		Total: sm.Total,
		Used:  sm.Used,
		Free:  sm.Free,
	}
	result.Percent.Swap = sm.UsedPercent

	diskUsage, err := disk.Usage("/")
	if err != nil {
		return result, err
	}
	result.Disk.UsageStat = UsageStat{
		Total: diskUsage.Total,
		Used:  diskUsage.Used,
		Free:  diskUsage.Free,
	}
	result.Percent.Disk = diskUsage.UsedPercent
	result.Disk.Read = backgroundStat.DiskReadSpeed
	result.Disk.Write = backgroundStat.DiskWriteSpeed

	result.Network.Rx = backgroundStat.NetRx
	result.Network.Tx = backgroundStat.NetTx
	result.Network.In = backgroundStat.NetTotalIn
	result.Network.Out = backgroundStat.NetTotalOut
	result.Network.Min = backgroundStat.NetMonthIn
	result.Network.Mout = backgroundStat.NetMonthOut

	hostinfo, err := host.Info()
	if err != nil {
		return result, err
	}
	result.Hostname = hostinfo.Hostname
	result.Uptime = hostinfo.Uptime
	return result, nil
}
