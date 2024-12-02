package charts

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"graph-viewer/logger"

	"github.com/go-echarts/go-echarts/v2/opts"
)

// Checks if all Y-axis values are strings
func areAllYValuesStrings(yAxis []opts.BarData) bool {
	for _, v := range yAxis {
		if _, ok := v.Value.(string); !ok {
			return false
		}
	}
	return true
}

func ContainsNonNumericData(data [][]string) bool {
	for _, row := range data[1:] { // Skip header
		if len(row) > 1 {
			if _, err := strconv.ParseFloat(row[1], 64); err != nil {
				return true // Found a non-numeric value
			}
		}
	}
	return false
}

// GenerateGraph creates a graph based on the selected type
func GenerateGraph(data [][]string, graphType string) (string, error) {
	switch graphType {
	case "Bar":
		return GenerateBarChart(data)
	case "Heatmap":
		return GenerateHeatmap(data)
	case "Kline":
		return GenerateKlineChart(data)
	case "Pie":
		return GeneratePieChart(data)
	case "Sankey":
		return GenerateSankeyChart(data)
	case "Overlap":
		return GenerateOverlapChart(data)
	case "Scatter3D":
		return GenerateScatter3D(data)
	case "Bar3D":
		return GenerateBar3DChart(data)
	case "ThemeRiver":
		return GenerateThemeRiverChart(data)
	default:
		return "", fmt.Errorf("unsupported graph type: %s", graphType)
	}
}

func generateBarChart(data [][]string) (string, error) {
	// Generate a bar chart (implement logic similar to your existing bar chart code)
	filePath := "bar_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Add bar chart rendering logic here...

	return filePath, nil
}

// Opens the chart.html file in the default browser
func ShowChartInBrowser(filePath string) error {
	logger.LogWithTrace(fmt.Sprintf("Attempting to open chart at: %s", filePath))

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Use "start" for Windows
		cmd = exec.Command("cmd", "/c", "start", filePath)
	case "darwin":
		// Use "open" for macOS
		cmd = exec.Command("open", filePath)
	case "linux":
		// Use "xdg-open" for Linux
		cmd = exec.Command("xdg-open", filePath)
	default:
		err := fmt.Errorf("unsupported platform: %s", runtime.GOOS)
		logger.LogErrorWithTrace(err)
		return err
	}

	// Start the command and handle errors
	err := cmd.Start()
	if err != nil {
		logger.LogErrorWithTrace(fmt.Errorf("Failed to open browser: %v", err))
		return err
	}

	logger.LogWithTrace("Successfully opened the chart in the browser.")
	return nil
}
