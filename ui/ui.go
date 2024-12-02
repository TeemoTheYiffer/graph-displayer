package ui

import (
	"errors"
	"fmt"
	"graph-viewer/charts"
	"graph-viewer/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SetupUI creates the main application UI
func SetupUI(window fyne.Window) {
	// Verify embedded files at startup
	verifyEmbeddedFiles()

	label := widget.NewLabel("Upload a CSV/XLS file to display an interactive graph.")
	fileButton := widget.NewButton("Select File", createFileHandler(window))

	content := container.NewVBox(label, fileButton)
	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 400))
}

// createFileHandler returns a function to handle file selection
func createFileHandler(window fyne.Window) func() {
	return func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				logger.LogErrorWithTrace(fmt.Errorf("file selection error: %v", err))
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return // User canceled
			}

			filePath := reader.URI().Path()
			defer reader.Close()

			headers, rows, err := readData(filePath)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			ShowHeaderSelection(headers, window, func(graphType string, xAxis string, yAxis string, zAxis string, limits map[string]int) {
				handleGraphGeneration(window, graphType, xAxis, yAxis, zAxis, headers, rows, limits)
			})
		}, window)
	}
}

// handleGraphGeneration processes the selected data and generates the graph
func handleGraphGeneration(
	window fyne.Window,
	graphType, xAxis, yAxis, zAxis string,
	headers []string,
	rows [][]string,
	limits map[string]int,
) {
	logger.LogWithTrace(fmt.Sprintf("Graph Type: %s, X-Axis: %s, Y-Axis: %s, Z-Axis: %s, Limits: %v",
		graphType, xAxis, yAxis, zAxis, limits))

	var (
		selectedData [][]string
		err          error
	)

	if graphType == "Scatter3D" || graphType == "Bar3D" {
		selectedData, err = extractSelectedData3D(xAxis, yAxis, zAxis, headers, rows, limits)
	} else {
		selectedData, err = extractSelectedData(xAxis, yAxis, headers, rows, limits)
	}

	if err != nil {
		logger.LogErrorWithTrace(fmt.Errorf("error extracting selected data: %v", err))
		dialog.ShowError(err, window)
		return
	}

	graphFile, err := charts.GenerateGraph(selectedData, graphType)
	if err != nil {
		logger.LogErrorWithTrace(fmt.Errorf("error generating graph: %v", err))
		dialog.ShowError(err, window)
		return
	}

	if err := charts.ShowChartInBrowser(graphFile); err != nil {
		logger.LogErrorWithTrace(fmt.Errorf("error displaying chart in browser: %v", err))
		dialog.ShowError(errors.New("failed to open chart in browser. Check logs for details."), window)
	}
}

// verifyEmbeddedFiles checks if embedded files are present
func verifyEmbeddedFiles() {
	entries, err := embeddedExamples.ReadDir("examples")
	if err != nil {
		logger.LogErrorWithTrace(fmt.Errorf("failed to read embedded examples directory: %v", err))
		return
	}

	logger.LogWithTrace(fmt.Sprintf("Found %d embedded files:", len(entries)))
	for _, entry := range entries {
		logger.LogWithTrace(fmt.Sprintf("- %s", entry.Name()))
	}

	expectedGifs := []string{
		"Bar.gif", "Heatmap.gif", "Scatter3D.gif", "Bar3D.gif",
		"Pie.gif", "Sankey.gif", "ThemeRiver.gif",
	}

	for _, gifName := range expectedGifs {
		if _, err := embeddedExamples.ReadFile(fmt.Sprintf("examples/%s", gifName)); err != nil {
			logger.LogErrorWithTrace(fmt.Errorf("failed to read embedded file %s: %v", gifName, err))
		}
	}
}
