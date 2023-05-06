package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func printStats() {
	// Get CPU usage
	cpuPercent, _ := cpu.Percent(time.Second, false)

	// Get memory usage
	memInfo, _ := mem.VirtualMemory()

	// Get disk usage
	diskInfo, _ := disk.Usage("/")

	// Print stats
	fmt.Printf("CPU usage: %.2f%%\n", cpuPercent[0])
	fmt.Printf("Memory usage: %.2f%%\n", memInfo.UsedPercent)
	fmt.Printf("Disk usage: %.2f%%\n", diskInfo.UsedPercent)
}

func main() {
	for {
		printStats()
		time.Sleep(5 * time.Second)
	}
}
