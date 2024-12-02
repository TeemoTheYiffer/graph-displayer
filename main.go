package main

import (
	"graph-viewer/ui"
	"os"
	"runtime/pprof"
	"time"

	"fyne.io/fyne/v2/app"
)

func main() {
	// CPU profiling
	cpuFile, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	// Memory profiling
	memFile, _ := os.Create("mem.prof")
	defer func() {
		pprof.WriteHeapProfile(memFile)
		memFile.Close()
	}()

	// Record startup time
	start := time.Now()

	application := app.NewWithID("com.git.teemotheyiffer.graph-viewer")
	window := application.NewWindow("Graph Viewer")

	ui.SetupUI(window)

	// Log startup time
	println("Startup time:", time.Since(start).Milliseconds(), "ms")

	window.ShowAndRun()
}
