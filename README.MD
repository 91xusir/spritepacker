# 🧩 SpritePacker

**SpritePacker** is a 2D sprite sheet packing tool designed for game developers.

------

## 🛠 Usage

### CLI Command

```bash
go install github.com/91xusir/spritepacker@latest
```
```bash
spritepacker -help
```
### 🛠️ Packing Options

| Parameter    | Description                                                                                                                           |
|--------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `-input`     | Input directory containing sprite images (required)                                                                                   |
| `-output`    | Output directory (default: `output`)                                                                                                  |
| `-maxw`      | Maximum atlas width (default: `2048`)                                                                                                 |
| `-maxh`      | Maximum atlas height (default: `2048`)                                                                                                |
| `-padding`   | Padding between sprites (default: `2`)                                                                                                |
| `-rotate`    | Enable sprite rotation for better space utilization                                                                                   |
| `-pot`       | Force power-of-two dimensions for atlas                                                                                               |
| `-name`      | Base name for output atlas (default: `atlas`)                                                                                         |
| `-sort`      | Sort images before packing (enabled by default)                                                                                       |
| `-trim`      | Trim transparent edges (enabled by default)                                                                                           |
| `-tolerance` | Alpha tolerance for trimming (0-255, default: `1`)                                                                                    |
| `-same`      | Enable duplicate image detection (enabled by default)                                                                                 |
| `-algo`      | Packing algorithm: 0 = Basic, 1 = Skyline, 2 = MaxRects                                                                               |
| `-heuristic` | MaxRects heuristic (only for `-algo 2`): 0 = BestShortSideFit, 1 = BestLongSideFit, 2 = BestAreaFit, 3 = BottomLeft, 4 = ContactPoint |

### 🛠️ Unpacking Options

| Parameter | Description                                        |
|-----------|----------------------------------------------------|
| `-unpack` | Path to atlas JSON file (required)                 |
| `-img`    | Path to atlas image (optional, inferred from JSON) |
| `-output` | Output directory (optional, inferred from JSON)    |

### 📦 Examples

```bash
# Pack all images in ./sprites with default settings
spritepacker -input ./sprites -output ./atlases

# Custom atlas size, padding, and enable rotation
spritepacker -input ./sprites -maxw 1024 -maxh 1024 -trim -

# Unpack atlas using atlas.json
spritepacker -unpack ./atlases/atlas.json -output ./unpacked
```

------


## 🧪 API

```bash
go get -u github.com/91xusir/spritepacker
```

Example Usage:

```go
package main

import (
	"github.com/91xusir/spritepacker/pack"
	"encoding/json"
	"os"
	"path/filepath"
)

func Pack() {
	// Step 1: Create packing options and set the maximum atlas size to 4096x4096
	options, _ := (pack.NewOptions().
		MaxSize(4096, 4096)).
	    Trim(true).
		Sort(true).
		Padding(-1).
		//AllowRotate(true).
		//Algorithm(pack.AlgoSkyline).
		//Heuristic(pack.ContactPointFit).
		//AutoSize(true).
		//PowerOfTwo(true).
		Validate()	// Optional: Validate the options. If invalid, an error will be returned.

	// Step 2: Create the sprite packer using the configured options
	spritePacker := pack.NewPacker(options)

	// Step 3: Collect all sprite image paths
	spriteImgPaths, _ := pack.GetFilesInDirectory("input")

	// Step 4: Pack all sprite images and generate atlas metadata and images
	spriteAtlasInfo, atlasImages, _ := spritePacker.PackSprites(spriteImgPaths)

	// Step 5: Save each atlas image to the "output" directory
	for i := range atlasImages {
		outputPath := filepath.Join("output", spriteAtlasInfo.Atlases[i].Name)
		_ = pack.SaveImg(outputPath, atlasImages[i], pack.PNG, pack.WithCLV(pack.DefaultCompression))
	}
	// Step 6: Save the atlas metadata as a JSON file
	jsonBytes, _ := json.MarshalIndent(spriteAtlasInfo, "", "  ")
	_ = os.WriteFile("output/atlas.json", jsonBytes, os.ModePerm)

}

func Unpack() {
	// default output and atlas images path is the same as the atlas.json path
	// if you want to change this, you can use the following code:
	// pack.UnpackSprites("output/atlas.json", pack.WithAtlasImgPath("output"), pack.WithOutputPath("output"))
	err := pack.UnpackSprites("output/atlas.json")
	if err != nil {
		return
	}
}
```

# ✅ TODO List

- [ ] Custom output format using Go templates
- [ ] more format support
- [ ] GUI 
