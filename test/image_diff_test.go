package spritepacker

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

// image diff test
func TestImageDiff(t *testing.T) {
	inputFolder := "input"
	outputFolder := "output"
	diffs := compareFolders(inputFolder, outputFolder)
	if len(diffs) > 0 {
		fmt.Println("The following images are different:")
		for _, name := range diffs {
			fmt.Println(" -", name)
		}
	} else {
		fmt.Println("All images are identical ✅")
	}
}

func isImageDifferent(img1Path, img2Path string) bool {
	file1, err := os.Open(img1Path)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", img1Path, err)
		return true
	}
	defer file1.Close()

	file2, err := os.Open(img2Path)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", img2Path, err)
		return true
	}
	defer file2.Close()

	img1, _, err := image.Decode(file1)
	if err != nil {
		fmt.Printf("Failed to decode %s: %v\n", img1Path, err)
		return true
	}

	img2, _, err := image.Decode(file2)
	if err != nil {
		fmt.Printf("Failed to decode %s: %v\n", img2Path, err)
		return true
	}

	if img1.Bounds().Size() != img2.Bounds().Size() {
		return true
	}

	bounds := img1.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img1.At(x, y) != img2.At(x, y) {
				return true
			}
		}
	}

	return false
}

func compareFolders(inputDir, outputDir string) []string {
	var differentFiles []string
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Println("Error reading input dir:", err)
		return nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		inputPath := filepath.Join(inputDir, entry.Name())
		outputPath := filepath.Join(outputDir, entry.Name())
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			fmt.Printf("Missing file: %s\n", outputPath)
			continue
		}
		if isImageDifferent(inputPath, outputPath) {
			differentFiles = append(differentFiles, entry.Name())
		}
	}
	return differentFiles
}
