package spritepacker

import (
	"github.com/91xusir/spritepacker/export"
	"github.com/91xusir/spritepacker/pack"
	"github.com/91xusir/spritepacker/utils"
	"path/filepath"
	"testing"
)

func TestSpritePacker(t *testing.T) {
	options := pack.NewOptions().MaxSize(4096, 4096).
		Trim(true).
		Sort(true).
		Padding(0).
		AllowRotate(false).
		Algorithm(pack.AlgoBasic).
		Heuristic(pack.ContactPointFit).
		ImgExt("webp").
		SameDetect(true).
		AutoSize(true).
		PowerOfTwo(false)

	atlasInfo, atlasImages, _ := pack.NewPacker(options).PackSprites("./input")

	for i := range atlasImages {
		outputPath := filepath.Join("output", atlasInfo.Atlases[i].Name)
		_ = utils.SaveImgByExt(outputPath, atlasImages[i], utils.WithCLV(utils.DefaultCompression))
	}

	_ = export.NewExportManager().Init().Export("output/atlas.tpsheet", atlasInfo)

}

func TestSpriteUnpack(t *testing.T) {
	// pack.UnpackSprites("output/atlas.json", pack.WithImgInput("output"), pack.WithOutput("output"))
	err := pack.UnpackSprites("output/atlas.tpsheet")
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
