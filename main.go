package main

import (
	"strings"
	"time"
	"fmt"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func getListeningPorts() string {
    conns, err := net.Connections("tcp")
    if err != nil {
        return fmt.Sprintf("Error fetching listening ports: %v", err)
    }

    var listeningPorts []string
    for _, conn := range conns {
        if conn.Status == "LISTEN" {
            listeningPorts = append(listeningPorts, fmt.Sprintf("%d", conn.Laddr.Port))
        }
    }
    return strings.Join(listeningPorts, ", ")
}


func printStats() (int, int, int) {
	// Get CPU usage
	cpuPercent, _ := cpu.Percent(time.Second, false)

	// Get memory usage
	memInfo, _ := mem.VirtualMemory()

	// Get disk usage
	diskInfo, _ := disk.Usage("/")

	// Convert CPU, memory, and disk usage to integers
	return int(cpuPercent[0]), int(memInfo.UsedPercent), int(diskInfo.UsedPercent)
}


func getCPUTemperature() float64 {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return 0.0
	}

	cpuTempSensorKeys := []string{
		"coretemp_packageid0",
		"coretemp",
		"k10temp",
		"Tdie",
	}

	for _, sensor := range sensors {
		for _, key := range cpuTempSensorKeys {
			if strings.Contains(sensor.SensorKey, key) {
				return sensor.Temperature
			}
		}
	}

	return 0.0
}

func main() {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	cpuGauge := widgets.NewGauge()
	cpuGauge.Title = " CPU Usage "
	cpuGauge.SetRect(0, 3, 80, 8)
	cpuGauge.Percent = 0
	cpuGauge.BarColor = termui.ColorBlue
	cpuGauge.BorderStyle.Fg = termui.ColorGreen
	cpuGauge.TitleStyle.Fg = termui.ColorWhite

	memGauge := widgets.NewGauge()
	memGauge.Title = " Memory Usage "
	memGauge.SetRect(0, 8, 80, 13)
	memGauge.Percent = 0
	memGauge.BarColor = termui.ColorBlue
	memGauge.BorderStyle.Fg = termui.ColorGreen
	memGauge.TitleStyle.Fg = termui.ColorWhite

	diskGauge := widgets.NewGauge()
	diskGauge.Title = " Disk Usage "
	diskGauge.SetRect(0, 13, 80, 18)
	diskGauge.Percent = 0
	diskGauge.BarColor = termui.ColorBlue
	diskGauge.BorderStyle.Fg = termui.ColorGreen
	diskGauge.TitleStyle.Fg = termui.ColorWhite

	cpuTempChart := widgets.NewPlot()
	cpuTempChart.Title = " CPU Temperature "
	cpuTempChart.SetRect(0, 23, 80, 40)
	cpuTempChart.BorderStyle.Fg = termui.ColorGreen
	cpuTempChart.TitleStyle.Fg = termui.ColorWhite
	cpuTempChart.LineColors[0] = termui.ColorBlue
	cpuTempChart.AxesColor = termui.ColorGreen
	cpuTempChart.Marker = widgets.MarkerDot
	cpuTempChart.Data = [][]float64{[]float64{}}
	

	title := widgets.NewParagraph()
	title.Text = " Sys Lizard "
	title.SetRect(0, 0, 80, 3)
	title.BorderStyle.Fg = termui.ColorGreen
	title.TitleStyle.Fg = termui.ColorGreen
	title.TextStyle.Fg = termui.ColorWhite
	title.Border = true

	listeningPortsText := widgets.NewParagraph()
	listeningPortsText.Title = " Listening Ports "
	listeningPortsText.SetRect(0, 18, 80, 23)
	listeningPortsText.BorderStyle.Fg = termui.ColorGreen
	listeningPortsText.TitleStyle.Fg = termui.ColorWhite
	listeningPortsText.TextStyle.Fg = termui.ColorWhite
	listeningPortsText.Border = true


	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(5 * time.Second)


	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}

		case <-ticker.C:
			cpuUsage, memUsage, diskUsage := printStats()
			cpuGauge.Percent = cpuUsage
			memGauge.Percent = memUsage
			diskGauge.Percent = diskUsage
			listeningPorts := getListeningPorts()
    		listeningPortsText.Text = listeningPorts

			cpuTemp := getCPUTemperature()
			cpuTempChart.Data[0] = append(cpuTempChart.Data[0], cpuTemp)
			if len(cpuTempChart.Data[0]) > 10 {
				cpuTempChart.Data[0] = cpuTempChart.Data[0][1:]
			}

			termui.Render(title, cpuGauge, memGauge, diskGauge, cpuTempChart, listeningPortsText)

		}
	}
}