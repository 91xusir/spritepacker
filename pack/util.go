package pack

import (
	"image"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func Parallel(start, end int, fn func(i int)) {
	numGoroutines := runtime.NumCPU()
	if end-start < numGoroutines {
		for i := start; i < end; i++ {
			fn(i)
		}
		return
	}
	var wg sync.WaitGroup
	batchSize := (end - start) / numGoroutines
	if batchSize < 1 {
		batchSize = 1
	}
	for i := start; i < end; i += batchSize {
		wg.Add(1)
		go func(from, to int) {
			defer wg.Done()
			for j := from; j < to && j < end; j++ {
				fn(j)
			}
		}(i, i+batchSize)
	}
	wg.Wait()
}

// GetOpaqueBounds returns the bounds of the opaque area of the image.
// The image is assumed to be in RGBA format.
func GetOpaqueBounds(img image.Image, tolerance uint8) image.Rectangle {
	bounds := img.Bounds()
	if bounds.Empty() {
		return bounds
	}

	// Initialize to opposite extremes
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y
	found := false

	// Helper function to update bounds
	updateBounds := func(x, y int) {
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		found = true
	}

	switch src := img.(type) {
	case *image.RGBA, *image.NRGBA:
		// Unified handling for RGBA and NRGBA
		var pix []uint8
		var stride int
		switch t := src.(type) {
		case *image.RGBA:
			pix = t.Pix
			stride = t.Stride
		case *image.NRGBA:
			pix = t.Pix
			stride = t.Stride
		}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			i := (y-bounds.Min.Y)*stride + (bounds.Min.X-bounds.Min.X)*4
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if pix[i+3] > tolerance {
					updateBounds(x, y)
				}
				i += 4
			}
		}

	default:
		// Generic case for other image types
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				_, _, _, a := img.At(x, y).RGBA()
				if uint8(a>>8) > tolerance { // Convert 16-bit alpha to 8-bit
					updateBounds(x, y)
				}
			}
		}
	}

	if !found {
		return bounds
	}
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

// GetFilesInDirectory reads files in specified directory and returns file names slice
func GetFilesInDirectory(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("Failed to read directory %s: %v", dirPath, err)
		return nil, err
	}
	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	return fileNames, nil
}

// GetLastFolderName 获取路径中的最后一个文件夹名称
func GetLastFolderName(path string) string {
	path = filepath.ToSlash(path)
	path = strings.TrimRight(path, "/")
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return "atlas"
}
