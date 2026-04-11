// Package sysinfo collects CPU, memory, disk, and uptime information.
// Used by the /status Telegram command and the heartbeat.
package sysinfo

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// Info holds a snapshot of system resource usage.
type Info struct {
	CPUPercent float64
	MemPercent float64
	MemUsed    uint64
	MemTotal   uint64
	Disks      []DiskUsage
	Uptime     time.Duration
	Platform   string
	OS         string
}

// DiskUsage holds disk usage for one mount point.
type DiskUsage struct {
	Path  string
	Free  uint64
	Total uint64
}

// Collect gathers the current system snapshot.
// On error for any sub-metric, it returns zero values for that metric
// rather than propagating the error — the caller always gets a partial result.
func Collect() Info {
	info := Info{
		Platform: runtime.GOOS + "/" + runtime.GOARCH,
	}

	// CPU usage (1-second average)
	if percents, err := cpu.Percent(1*time.Second, false); err == nil && len(percents) > 0 {
		info.CPUPercent = percents[0]
	}

	// Memory
	if vm, err := mem.VirtualMemory(); err == nil {
		info.MemPercent = vm.UsedPercent
		info.MemUsed    = vm.Used
		info.MemTotal   = vm.Total
	}

	// Disk — all mounted partitions
	info.Disks = collectDisks()

	// Uptime
	if uptimeSec, err := host.Uptime(); err == nil {
		info.Uptime = time.Duration(uptimeSec) * time.Second
	}

	// OS
	if hi, err := host.Info(); err == nil {
		info.OS = hi.Platform + " " + hi.PlatformVersion
	}

	return info
}

func collectDisks() []DiskUsage {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var disks []DiskUsage

	for _, p := range partitions {
		// Skip duplicates and pseudo-filesystems
		if seen[p.Mountpoint] {
			continue
		}
		if isSkippedFS(p.Fstype) {
			continue
		}
		seen[p.Mountpoint] = true

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}

		disks = append(disks, DiskUsage{
			Path:  p.Mountpoint,
			Free:  usage.Free,
			Total: usage.Total,
		})
	}

	// Limit to 3 most relevant disks to keep the message concise
	if len(disks) > 3 {
		disks = disks[:3]
	}
	return disks
}

func isSkippedFS(fstype string) bool {
	skip := map[string]bool{
		"proc": true, "sysfs": true, "devtmpfs": true, "devpts": true,
		"tmpfs": true, "cgroup": true, "cgroup2": true, "pstore": true,
		"efivarfs": true, "bpf": true, "securityfs": true, "hugetlbfs": true,
		"mqueue": true, "debugfs": true, "tracefs": true, "fusectl": true,
		"configfs": true, "squashfs": true, "overlay": true,
	}
	return skip[fstype]
}
