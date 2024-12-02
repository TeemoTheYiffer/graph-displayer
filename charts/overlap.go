package charts

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateOverlapChart creates a chart with overlapping series (e.g., Bar and Line)
func GenerateOverlapChart(data [][]string) (string, error) {
	if len(data) < 2 || len(data[0]) < 3 {
		return "", fmt.Errorf("overlap chart requires at least 3 columns: X, Y1, Y2")
	}

	// Extract data for the charts
	xLabels := []string{}
	barValues := []opts.BarData{}
	lineValues := []opts.LineData{}

	for i, row := range data[1:] { // Skip header row
		if len(row) < 3 {
			continue // Skip rows with insufficient columns
		}

		xLabels = append(xLabels, row[0]) // First column as X-axis label

		// Parse values for bar and line
		y1, err1 := parseNumericValue(row[1])
		y2, err2 := parseNumericValue(row[2])

		if err1 != nil || err2 != nil {
			fmt.Printf("Skipping invalid row %d: %v, %v\n", i+1, err1, err2)
			continue
		}

		barValues = append(barValues, opts.BarData{Value: y1})
		lineValues = append(lineValues, opts.LineData{Value: y2})
	}

	if len(xLabels) == 0 || len(barValues) == 0 || len(lineValues) == 0 {
		return "", fmt.Errorf("no valid data for overlap chart")
	}

	// Create Bar chart
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Overlap Chart",
			Subtitle: "Bar and Line",
		}),
		charts.WithXAxisOpts(opts.XAxis{Name: "X-axis"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Y-axis"}),
	)
	bar.SetXAxis(xLabels).AddSeries("Bar", barValues)

	// Create Line chart
	line := charts.NewLine()
	line.SetXAxis(xLabels).AddSeries("Line", lineValues)

	// Overlap the charts
	bar.Overlap(line)

	// Render to file
	filePath := "overlap_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = bar.Render(file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
