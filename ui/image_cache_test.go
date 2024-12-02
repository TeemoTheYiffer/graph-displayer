// image_cache_test.go
package ui

import (
	"fmt"
	"testing"
)

func TestEmbeddedFilesPresent(t *testing.T) {
	// List of expected Images
	expectedPNGs := []string{
		"Bar.png",
		"Heatmap.png",
		"Scatter3D.png",
		"Bar3D.png",
		"Pie.png",
		"Sankey.png",
		"ThemeRiver.png",
	}

	// Get all embedded files
	entries, err := embeddedExamples.ReadDir("examples")
	if err != nil {
		t.Fatalf("Failed to read embedded directory: %v", err)
	}

	fmt.Printf("Total files found in embedded filesystem: %d\n", len(entries))

	// Create a map of found files
	foundFiles := make(map[string]bool)
	for _, entry := range entries {
		name := entry.Name()
		foundFiles[name] = true

		// Try to read the file data
		data, err := embeddedExamples.ReadFile("examples/" + name)
		if err != nil {
			t.Errorf("Failed to read %s: %v", name, err)
		} else {
			fmt.Printf("Successfully read %s (%d bytes)\n", name, len(data))
		}
	}

	// Check for missing files
	var missingFiles []string
	for _, expected := range expectedPNGs {
		if !foundFiles[expected] {
			missingFiles = append(missingFiles, expected)
		}
	}

	if len(missingFiles) > 0 {
		t.Errorf("Missing expected PNG files: %v", missingFiles)
	}
}
