package main

import (
	"fmt"
	"time"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
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

func printCPUTemperature() {
    sensors, err := host.SensorsTemperatures()
    if err != nil {
        fmt.Println("Error fetching CPU temperature:", err)
        return
    }

    cpuTempSensorKeys := []string{
        "coretemp_packageid0",
        "coretemp",
        "k10temp",
        "Tdie",
    }

    found := false
    for _, sensor := range sensors {
        for _, key := range cpuTempSensorKeys {
            if strings.Contains(sensor.SensorKey, key) {
                fmt.Printf("CPU temperature: %.1f°C\n", sensor.Temperature)
                found = true
                break
            }
        }
        if found {
            break
        }
    }

    if !found {
        fmt.Println("CPU temperature sensor not found. Please check the sensor key.")
    }
}




func main() {
	for {
		printStats()
		printCPUTemperature()
		time.Sleep(5 * time.Second)
	}
}
