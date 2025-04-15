package spritepacker

import (
	"encoding/json"
	"github.com/91xusir/spritepacker/pack"
	"github.com/91xusir/spritepacker/utils"
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
		AllowRotate(false).
		Algorithm(pack.AlgoBasic).
		Heuristic(pack.ContactPointFit).
		SameDetect(true).
		AutoSize(true).
		PowerOfTwo(false)

	spriteAtlasInfo, atlasImages, _ := pack.NewPacker(options).PackSprites("./input")

	for i := range atlasImages {
		outputPath := filepath.Join("output", spriteAtlasInfo.Atlases[i].Name)
		_ = utils.SaveImgByExt(outputPath, atlasImages[i], utils.WithCLV(utils.DefaultCompression))
	}

	jsonBytes, _ := json.MarshalIndent(spriteAtlasInfo, "", "  ")
	_ = os.WriteFile("output/atlas.json", jsonBytes, 0644)

}

func TestSpriteUnpack(t *testing.T) {
	// pack.UnpackSprites("output/atlas.json", pack.WithImg("output"), pack.WithOutput("output"))
	err := pack.UnpackSprites("output/atlas.json")
	if err != nil {
		t.Errorf("Failed to unpack sprites: %v", err)
		return
	}
}

func TestImageDiff(t *testing.T) {
	inputFolder := "input"
	outputFolder := "output"
	diffs := utils.CompareImgFormFolders(inputFolder, outputFolder)
	if len(diffs) > 0 {
		t.Errorf("Found %d different images:\n", len(diffs))
	} else {
		t.Logf("All images are the same.\n")
	}
}
