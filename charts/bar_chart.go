package charts

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateBarChart creates an HTML bar chart from the given data
func GenerateBarChart(data [][]string) (string, error) {
	// Validate input data
	if len(data) < 2 {
		return "", fmt.Errorf("insufficient data for bar chart")
	}

	// Extract X-axis labels and Y-axis values
	xLabels := []string{}
	yValues := []opts.BarData{}

	for _, row := range data[1:] { // Skip the header row
		if len(row) < 2 {
			continue // Skip rows without enough columns
		}

		xLabels = append(xLabels, row[0]) // First column as X-axis label

		// Parse Y-axis value (skip if invalid)
		yValue := row[1]
		value, err := parseNumericValue(yValue)
		if err != nil {
			fmt.Printf("Skipping invalid Y-axis value: %s\n", yValue)
			continue
		}
		yValues = append(yValues, opts.BarData{Value: value})
	}

	if len(xLabels) == 0 || len(yValues) == 0 {
		return "", fmt.Errorf("no valid data for bar chart")
	}

	// Create a new bar chart
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Bar Chart",
			Subtitle: "",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Categories",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Values",
		}),
	)

	// Add data to the chart
	bar.SetXAxis(xLabels).AddSeries("Data", yValues)

	// Render the chart to an HTML file
	filePath := "bar_chart.html"
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

// Converts a string to a float64, handling potential errors
func parseNumericValue(input string) (float64, error) {
	// Attempt to parse the string as a float
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value: %s", input)
	}
	return value, nil
}
