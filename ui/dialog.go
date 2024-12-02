package ui

import (
	"errors"
	"fmt"
	"graph-viewer/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// GraphTypeInfo holds metadata about different graph types
type GraphTypeInfo struct {
	name        string
	needs3D     bool
	description string
}

// Available graph types and their metadata
var graphTypeInfos = map[string]GraphTypeInfo{
	"Bar":        {"Bar Chart", false, "Simple bar chart for comparing categories"},
	"Heatmap":    {"Heat Map", false, "Visualize data density and patterns"},
	"Scatter3D":  {"3D Scatter Plot", true, "Three-dimensional scatter visualization"},
	"Bar3D":      {"3D Bar Chart", true, "Three-dimensional bar visualization"},
	"Pie":        {"Pie Chart", false, "Show proportion between categories"},
	"Sankey":     {"Sankey Diagram", false, "Visualize flow between categories"},
	"ThemeRiver": {"Theme River", false, "Show changes over time"},
}

// ShowHeaderSelection creates and shows the graph type and axis selection dialog
func ShowHeaderSelection(headers []string, window fyne.Window, callback func(graphType string, xAxis string, yAxis string, zAxis string, limits map[string]int)) {
	// Create UI components
	previewContainer := container.New(layout.NewHBoxLayout(),
		widget.NewLabel("Select a graph type to see preview"))
	descriptionLabel := widget.NewLabel("")
	axisSelectors := container.New(layout.NewVBoxLayout())

	// Create selectors
	graphTypeSelector := createGraphTypeSelector()
	xAxisSelector := widget.NewSelect(headers, nil)
	yAxisSelector := widget.NewSelect(headers, nil)
	zAxisSelector := widget.NewSelect(headers, nil)

	// Update form based on graph type selection
	updateForm := func(graphType string) {
		if graphType == "" {
			return
		}

		graphInfo := graphTypeInfos[graphType]
		descriptionLabel.SetText(graphInfo.description)

		updateAxisSelectors(axisSelectors, xAxisSelector, yAxisSelector, zAxisSelector, graphInfo.needs3D)
		updatePreviewImage(graphType, previewContainer)
	}

	// Set up callbacks
	graphTypeSelector.OnChanged = updateForm

	// Create dialog layout
	form := container.New(layout.NewVBoxLayout(),
		widget.NewForm(widget.NewFormItem("Graph Type", graphTypeSelector)),
		descriptionLabel,
		axisSelectors,
		previewContainer,
	)

	// Create and show dialog
	showSelectionDialog(window, form, graphTypeSelector, xAxisSelector, yAxisSelector, zAxisSelector, callback)
}

func updatePreviewImage(graphType string, container *fyne.Container) {
	img, err := loadImageFromCache(graphType)
	if err != nil {
		logger.LogErrorWithTrace(fmt.Errorf("failed to update preview image: %w", err))
		container.RemoveAll()
		container.Add(widget.NewLabel(fmt.Sprintf("No preview available for %s", graphType)))
		container.Refresh()
		return
	}

	container.RemoveAll()
	container.Add(img)
	container.Refresh()
}

// Helper functions
func createGraphTypeSelector() *widget.Select {
	var types []string
	for gType := range graphTypeInfos {
		types = append(types, gType)
	}
	selector := widget.NewSelect(types, nil)
	selector.SetSelected("Bar")
	return selector
}

func updateAxisSelectors(container *fyne.Container, x, y, z *widget.Select, needs3D bool) {
	container.RemoveAll()
	container.Add(widget.NewForm(
		&widget.FormItem{Text: "X Axis", Widget: x},
		&widget.FormItem{Text: "Y Axis", Widget: y},
	))
	if needs3D {
		container.Add(widget.NewForm(
			&widget.FormItem{Text: "Z Axis", Widget: z},
		))
	}
	container.Refresh()
}

func showSelectionDialog(
	window fyne.Window,
	content *fyne.Container,
	graphType, xAxis, yAxis, zAxis *widget.Select,
	callback func(string, string, string, string, map[string]int),
) {
	dialog := dialog.NewCustomConfirm(
		"Select Graph Type and Axes",
		"Create Graph",
		"Cancel",
		content,
		func(confirmed bool) {
			if !confirmed {
				return
			}

			if err := validateSelections(graphType, xAxis, yAxis, zAxis); err != nil {
				dialog.ShowError(err, window)
				return
			}

			callback(
				graphType.Selected,
				xAxis.Selected,
				yAxis.Selected,
				zAxis.Selected,
				nil,
			)
		},
		window,
	)

	dialog.Resize(fyne.NewSize(600, 600))
	dialog.Show()
}

func validateSelections(graphType, xAxis, yAxis, zAxis *widget.Select) error {
	if graphType.Selected == "" || xAxis.Selected == "" || yAxis.Selected == "" {
		return errors.New("please select graph type, X axis, and Y axis")
	}

	graphInfo := graphTypeInfos[graphType.Selected]
	if graphInfo.needs3D && zAxis.Selected == "" {
		return errors.New("please select Z axis for 3D visualization")
	}

	return nil
}
