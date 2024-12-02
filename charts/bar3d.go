package charts

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateBar3DChart creates a Bar3D chart from the given data
func GenerateBar3DChart(data [][]string) (string, error) {
	if len(data) < 2 || len(data[0]) < 3 {
		return "", fmt.Errorf("bar3D chart requires at least 3 columns: X, Y, Z")
	}

	// Extract data points
	points := []opts.Chart3DData{}
	for i, row := range data[1:] { // Skip header row
		if len(row) < 3 {
			continue // Skip rows with insufficient columns
		}

		x, errX := parseNumericValue(row[0])
		y, errY := parseNumericValue(row[1])
		z, errZ := parseNumericValue(row[2])

		if errX != nil || errY != nil || errZ != nil {
			fmt.Printf("Skipping invalid row %d: %v, %v, %v\n", i+1, errX, errY, errZ)
			continue
		}

		points = append(points, opts.Chart3DData{
			Value: toInterfaceSlice(x, y, z), // Convert to []interface{}
		})
	}

	if len(points) == 0 {
		return "", fmt.Errorf("no valid data for Bar3D chart")
	}

	// Create Bar3D chart
	bar3D := charts.NewBar3D()
	bar3D.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Bar3D Chart",
			Subtitle: "",
		}),
		charts.WithXAxis3DOpts(opts.XAxis3D{Name: "X"}),
		charts.WithYAxis3DOpts(opts.YAxis3D{Name: "Y"}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{Name: "Z"}),
	)

	// Add data to the chart
	bar3D.AddSeries("Bar3D", points)

	// Render to file
	filePath := "bar3d_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = bar3D.Render(file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
