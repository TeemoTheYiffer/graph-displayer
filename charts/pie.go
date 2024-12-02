package charts

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GeneratePieChart creates a Pie chart from the given data
func GeneratePieChart(data [][]string) (string, error) {
	if len(data) < 2 || len(data[0]) < 2 {
		return "", fmt.Errorf("pie chart requires at least 2 columns: Category, Value")
	}

	// Extract categories and values
	items := []opts.PieData{}
	for i, row := range data[1:] { // Skip header row
		if len(row) < 2 {
			continue // Skip rows with insufficient columns
		}

		category := row[0]
		value, err := parseNumericValue(row[1])
		if err != nil {
			fmt.Printf("Skipping invalid row %d: %v\n", i+1, err)
			continue
		}

		items = append(items, opts.PieData{Name: category, Value: value})
	}

	if len(items) == 0 {
		return "", fmt.Errorf("no valid data for Pie chart")
	}

	// Create Pie chart
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Pie Chart",
		}),
	)

	pie.AddSeries("Pie", items)

	// Render to file
	filePath := "pie_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = pie.Render(file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
