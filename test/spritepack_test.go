package spritepacker

import (
	"os"
	"path/filepath"
	"spritepacker/pack"
	"testing"
)

func TestSpritePacker(t *testing.T) {
	// first, set options
	options := pack.NewOptions().MaxSize(1024, 1024)
	// you can validate options to check if options are valid
	// if options are not valid, it will return error
	options, _ = options.Validate()

	// then, create sprite packer
	spritePacker := pack.NewPacker(options)

	// get all sprite image paths in folder
	spriteImgPaths, _ := pack.GetFilesInDirectory("input")

	spriteAtlasInfo, atlasImages, _ := spritePacker.PackSprites(spriteImgPaths)

	t.Logf("spriteAtilas: %v", spriteAtlasInfo)

	for i := range atlasImages {
		// you must use spriteAtlasInfo to get atlas name
		// that unpack function can use to get atlas image
		path := filepath.Join("output", spriteAtlasInfo.Atlases[i].Name)
		file, _ := os.Create(path)
		_ = pack.EncImg(file, atlasImages[i], pack.PNG, pack.WithCLV(pack.NoCompression))
	}
}
