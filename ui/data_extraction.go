package ui

import (
	"fmt"
	"strconv"
)

// validateNumeric checks if a string can be parsed as a float
func validateNumeric(value string) error {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("value '%s' is not a valid number", value)
	}
	return nil
}

// extractSelectedData3D extracts data for 3D visualizations
func extractSelectedData3D(xAxis, yAxis, zAxis string, headers []string, rows [][]string, limits map[string]int) ([][]string, error) {
	xIndex, yIndex, zIndex := -1, -1, -1

	// Find column indices
	for i, header := range headers {
		switch header {
		case xAxis:
			xIndex = i
		case yAxis:
			yIndex = i
		case zAxis:
			zIndex = i
		}
	}

	if xIndex == -1 || yIndex == -1 || zIndex == -1 {
		return nil, fmt.Errorf("invalid column selection for X, Y, or Z axis")
	}

	// Apply row limits
	maxRows := len(rows)
	if limit, ok := limits["X"]; ok && limit > 0 && limit < maxRows {
		maxRows = limit
	}

	selectedData := make([][]string, 0, maxRows)
	for i, row := range rows {
		if i >= maxRows {
			break
		}

		// Validate numeric data
		if err := validateNumeric(row[yIndex]); err != nil {
			return nil, fmt.Errorf("row %d: %v", i+1, err)
		}
		if err := validateNumeric(row[zIndex]); err != nil {
			return nil, fmt.Errorf("row %d: %v", i+1, err)
		}

		selectedData = append(selectedData, []string{row[xIndex], row[yIndex], row[zIndex]})
	}

	return selectedData, nil
}

// extractSelectedData extracts data for 2D visualizations
func extractSelectedData(xAxis, yAxis string, headers []string, rows [][]string, limits map[string]int) ([][]string, error) {
	xIndex, yIndex := -1, -1

	// Find column indices
	for i, header := range headers {
		switch header {
		case xAxis:
			xIndex = i
		case yAxis:
			yIndex = i
		}
	}

	if xIndex == -1 || yIndex == -1 {
		return nil, fmt.Errorf("invalid column selection for X or Y axis")
	}

	// Apply row limits
	maxRows := len(rows)
	if limit, ok := limits["X"]; ok && limit > 0 && limit < maxRows {
		maxRows = limit
	}

	selectedData := make([][]string, 0, maxRows)
	for i, row := range rows {
		if i >= maxRows {
			break
		}

		// Validate numeric data for Y axis
		if err := validateNumeric(row[yIndex]); err != nil {
			return nil, fmt.Errorf("row %d: %v", i+1, err)
		}

		selectedData = append(selectedData, []string{row[xIndex], row[yIndex]})
	}

	return selectedData, nil
}
