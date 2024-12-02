package ui

import (
	"embed"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

//go:embed examples/*.png
var embeddedExamples embed.FS

// imageCache implements a simple cache for image data
type imageCache struct {
	cache map[string][]byte
	mutex sync.RWMutex
}

var globalImageCache = &imageCache{
	cache: make(map[string][]byte),
}

// getImage loads an image from the cache or embedded filesystem
func (c *imageCache) getImage(name string) ([]byte, error) {
	c.mutex.RLock()
	if data, ok := c.cache[name]; ok {
		c.mutex.RUnlock()
		return data, nil
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if data, ok := c.cache[name]; ok {
		return data, nil
	}

	data, err := embeddedExamples.ReadFile(fmt.Sprintf("examples/%s.png", name))
	if err != nil {
		return nil, err
	}

	c.cache[name] = data
	return data, nil
}

// loadImageFromCache loads and creates an image resource
func loadImageFromCache(graphType string) (*canvas.Image, error) {
	imgData, err := globalImageCache.getImage(graphType)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}

	res := fyne.NewStaticResource(fmt.Sprintf("%s.png", graphType), imgData)
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(200, 150))

	return img, nil
}
