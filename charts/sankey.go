package charts

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateSankeyChart creates a Sankey chart from the given data
func GenerateSankeyChart(data [][]string) (string, error) {
	if len(data) < 2 || len(data[0]) < 3 {
		return "", fmt.Errorf("sankey chart requires at least 3 columns: Source, Target, Value")
	}

	// Create nodes and links
	nodesMap := map[string]struct{}{}
	links := []opts.SankeyLink{}

	for i, row := range data[1:] { // Skip header row
		if len(row) < 3 {
			continue // Skip rows with insufficient columns
		}

		source := row[0]
		target := row[1]
		value, err := parseNumericValue(row[2])
		if err != nil {
			fmt.Printf("Skipping invalid row %d: %v\n", i+1, err)
			continue
		}

		// Add source and target to nodes map
		nodesMap[source] = struct{}{}
		nodesMap[target] = struct{}{}

		// Add link
		links = append(links, opts.SankeyLink{
			Source: source,
			Target: target,
			Value:  float32(value), // Explicitly cast float64 to float32
		})
	}

	if len(nodesMap) == 0 || len(links) == 0 {
		return "", fmt.Errorf("no valid data for Sankey chart")
	}

	// Create nodes from map
	nodes := []opts.SankeyNode{}
	for node := range nodesMap {
		nodes = append(nodes, opts.SankeyNode{Name: node})
	}

	// Create Sankey chart
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Sankey Chart",
		}),
	)

	// Add nodes and links to the chart
	sankey.AddSeries("Sankey", nodes, links)

	// Render to file
	filePath := "sankey_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = sankey.Render(file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
