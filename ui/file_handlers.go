package ui

import (
	"fmt"
	"graph-viewer/logger"
	"os"
	"path/filepath"

	"encoding/csv"

	"github.com/xuri/excelize/v2"
)

// readData parses the file and returns data for charting
func readData(filePath string) ([]string, [][]string, error) {
	ext := filepath.Ext(filePath)
	var data [][]string
	var err error

	if ext == ".csv" {
		data, err = readCSV(filePath)
	} else if ext == ".xlsx" {
		data, err = readXLSX(filePath)
	} else {
		return nil, nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	if err != nil {
		return nil, nil, err
	}

	if len(data) < 2 {
		return nil, nil, fmt.Errorf("insufficient rows in file")
	}

	headers := data[0]
	rows := data[1:]

	if len(headers) == 0 {
		return nil, nil, fmt.Errorf("file contains no headers")
	}

	// Validate row lengths
	for i, row := range rows {
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("row %d has %d columns, expected %d", i+1, len(row), len(headers))
		}
	}

	return headers, rows, nil
}

// readCSV reads data from a CSV file
func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Check file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() > 100*1024*1024 {
		logger.LogWithTrace("Warning: Processing large file, this may take some time")
	}

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

// readXLSX reads data from an XLSX file
func readXLSX(filePath string) ([][]string, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file.GetRows(file.GetSheetName(file.GetActiveSheetIndex()))
}
