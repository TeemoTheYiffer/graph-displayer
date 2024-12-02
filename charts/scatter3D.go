package charts

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// toInterfaceSlice converts a slice of float64 to a slice of interface{}
func toInterfaceSlice(values ...float64) []interface{} {
	interfaceSlice := make([]interface{}, len(values))
	for i, v := range values {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

// GenerateScatter3D creates an HTML Scatter3D chart from the given data
func GenerateScatter3D(data [][]string) (string, error) {
	if len(data) < 2 || len(data[0]) < 3 {
		return "", fmt.Errorf("scatter3D requires at least 3 columns: X, Y, and Z")
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
		return "", fmt.Errorf("no valid data for Scatter3D")
	}

	// Create scatter3D chart
	scatter := charts.NewScatter3D()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Scatter3D",
			Subtitle: "",
		}),
		charts.WithXAxis3DOpts(opts.XAxis3D{Name: "X"}),
		charts.WithYAxis3DOpts(opts.YAxis3D{Name: "Y"}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{Name: "Z"}),
	)

	// Add data to the scatter3D chart
	scatter.AddSeries("Scatter3D", points)

	// Render the chart to an HTML file
	filePath := "scatter3d_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = scatter.Render(file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
