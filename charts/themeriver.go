package charts

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateThemeRiverChart creates a ThemeRiver chart from the given data
func GenerateThemeRiverChart(data [][]string) (string, error) {
	if len(data) < 2 || len(data[0]) < 3 {
		return "", fmt.Errorf("themeriver chart requires at least 3 columns: Time, Value, Category")
	}

	// Extract data points
	points := []opts.ThemeRiverData{}
	for i, row := range data[1:] { // Skip header row
		if len(row) < 3 {
			continue // Skip rows with insufficient columns
		}

		time := row[0]
		value, err := parseNumericValue(row[1])
		if err != nil {
			fmt.Printf("Skipping invalid row %d: %v\n", i+1, err)
			continue
		}
		category := row[2]

		points = append(points, opts.ThemeRiverData{
			Name:  category, // The category or type of the flow
			Value: value,    // Numeric value
			Date:  time,     // Time or date for the flow
		})
	}

	if len(points) == 0 {
		return "", fmt.Errorf("no valid data for ThemeRiver chart")
	}

	// Create ThemeRiver chart
	themeriver := charts.NewThemeRiver()
	themeriver.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "ThemeRiver Chart",
		}),
	)

	themeriver.AddSeries("ThemeRiver", points)

	// Render to file
	filePath := "themeriver_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = themeriver.Render(file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
