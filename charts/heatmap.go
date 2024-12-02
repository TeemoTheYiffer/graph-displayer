package charts

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateHeatmap creates an HTML heatmap chart from the given data
func GenerateHeatmap(data [][]string) (string, error) {
	if len(data) < 2 {
		return "", fmt.Errorf("insufficient data for heatmap")
	}

	// Extract headers and values
	headers := data[0]
	values := [][]opts.HeatMapData{}

	for i, row := range data[1:] { // Skip header row
		if len(row) != len(headers) {
			continue // Skip rows with mismatched column lengths
		}

		rowValues := []opts.HeatMapData{}
		for j, cell := range row {
			value, err := parseNumericValue(cell)
			if err != nil {
				fmt.Printf("Skipping invalid value at row %d, col %d: %s\n", i+1, j+1, cell)
				continue
			}
			rowValues = append(rowValues, opts.HeatMapData{
				Value: []interface{}{j, i, value}, // Column (x), Row (y), Value (z)
			})
		}
		values = append(values, rowValues)
	}

	if len(values) == 0 {
		return "", fmt.Errorf("no valid data for heatmap")
	}

	// Create heatmap chart
	heatmap := charts.NewHeatMap()
	heatmap.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Heatmap",
			Subtitle: "",
		}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Columns", Type: "category", Data: headers}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Rows", Type: "category"}),
	)

	// Flatten data for heatmap series
	flattenedValues := []opts.HeatMapData{}
	for _, row := range values {
		flattenedValues = append(flattenedValues, row...)
	}

	heatmap.AddSeries("Heatmap", flattenedValues)

	// Render the chart to an HTML file
	filePath := "heatmap_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = heatmap.Render(file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
