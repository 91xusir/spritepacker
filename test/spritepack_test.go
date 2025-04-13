package spritepacker

import (
	"encoding/json"
	"fmt"
	"github.com/91xusir/spritepacker/pack"
	"os"
	"path/filepath"
	"testing"
)

func TestSpritePacker(t *testing.T) {
	// Step 1: Create packing options and set the maximum atlas size to 4096x4096
	options := pack.NewOptions().MaxSize(4096, 4096).
		Trim(true).
		Sort(true).
		Padding(0).
		AllowRotate(true).
		Algorithm(pack.AlgoMaxRects).
		Heuristic(pack.BestAreaFit).
		AutoSize(true).
		PowerOfTwo(true)

	// Optional: Validate the options. If invalid, an error will be returned.
	// Here we ignore the error for demonstration purposes.
	options, _ = options.Validate()

	// Step 2: Create the sprite packer using the configured options
	spritePacker := pack.NewPacker(options)

	// Step 3: Collect all sprite image paths from the "input" directory
	spriteImgPaths, _ := pack.GetFilesInDirectory("input")

	// Step 4: Pack all sprite images and generate atlas metadata and images
	spriteAtlasInfo, atlasImages, _ := spritePacker.PackSprites(spriteImgPaths)

	t.Logf("spriteAtlas: %v", spriteAtlasInfo)

	// Step 5: Save each atlas image to the "output" directory
	for i := range atlasImages {
		// If you plan to unpack the atlas using spriteAtlasInfo later,
		// you should use the atlas name from spriteAtlasInfo as the filename.
		// Only the name is required, and the file extension will be determined
		// based on the format passed (e.g., pack.PNG, pack.JPEG, etc.).
		// WithCLV is function to set the compression level of the image.
		outputPath := filepath.Join("output", spriteAtlasInfo.Atlases[i].Name)
		_ = pack.SaveImg(outputPath, atlasImages[i], pack.PNG, pack.WithCLV(pack.DefaultCompression))
	}

	// Step 6: Save the atlas metadata as a JSON file
	// By default, all atlas information is saved into a single file.
	// You can read spriteAtlasInfo to save each atlas information separately.
	// Correspondingly, you would need to customize the UnpackSprites method
	// to read each atlas information and save each atlas image.
	jsonBytes, _ := json.MarshalIndent(spriteAtlasInfo, "", "  ")
	_ = os.WriteFile("output/atlas.json", jsonBytes, 0644)

}

func TestSpriteUnpack(t *testing.T) {
	// default output and atlas images path is the same as the atlas.json path
	// if you want to change this, you can use the following code:
	// pack.UnpackSprites("output/atlas.json", pack.WithAtlasImgPath("output"), pack.WithOutputPath("output"))
	err := pack.UnpackSprites("output/atlas.json")
	if err != nil {
		t.Errorf("Failed to unpack sprites: %v", err)
		return
	}
	// Then you can use TestImageDiff to compare the output images with the input images
	// to verify that the unpacking process is correct.
}

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

	img1, err := pack.DecImg(file1)
	if err != nil {
		fmt.Printf("Failed to decode %s: %v\n", img1Path, err)
		return true
	}

	img2, err := pack.DecImg(file2)
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
