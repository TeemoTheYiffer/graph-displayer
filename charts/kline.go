package charts

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateKlineChart creates a Kline chart from financial data
func GenerateKlineChart(data [][]string) (string, error) {
	if len(data) < 2 || len(data[0]) < 5 {
		return "", fmt.Errorf("kline chart requires at least 5 columns: Date, Open, Close, Low, High")
	}

	// Extract data
	values := []opts.KlineData{}
	xLabels := []string{}

	for i, row := range data[1:] { // Skip header row
		if len(row) < 5 {
			continue // Skip rows with insufficient columns
		}

		// Parse values
		date := row[0]
		open, err1 := parseNumericValue(row[1])
		close, err2 := parseNumericValue(row[2])
		low, err3 := parseNumericValue(row[3])
		high, err4 := parseNumericValue(row[4])

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Printf("Skipping invalid row %d: %v, %v, %v, %v\n", i+1, err1, err2, err3, err4)
			continue
		}

		xLabels = append(xLabels, date)
		values = append(values, opts.KlineData{Value: [4]float64{open, close, low, high}})
	}

	if len(values) == 0 {
		return "", fmt.Errorf("no valid data for kline chart")
	}

	// Create Kline chart
	kline := charts.NewKLine()
	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Kline Chart",
			Subtitle: "Financial Data",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Data: xLabels,
		}),
	)

	kline.AddSeries("Kline", values)

	// Render to file
	filePath := "kline_chart.html"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = kline.Render(file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
