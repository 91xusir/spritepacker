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

| Parameter | Type   | Description                                                                                                         |
|-----------|--------|---------------------------------------------------------------------------------------------------------------------|
| -i        | string | Input directory containing sprite images (required for packing)                                                     |
| -o        | string | Output directory (default "output")                                                                                 |
| -maxw     | int    | Maximum atlas width (default 2048)                                                                                  |
| -maxh     | int    | Maximum atlas height (default 2048)                                                                                 |
| -pad      | int    | Padding between sprites (default 0)                                                                                 |
| -auto     | bool   | Automatically adjust atlas size (default true)                                                                      |
| -rot      | bool   | Allow sprite rotation to save space (default false)                                                                 |
| -pot      | bool   | Force power-of-two atlas dimensions (default false)                                                                 |
| -name     | string | Base name for output files (default "atlas")                                                                        |
| -sort     | bool   | Sorts sprites before packing (default true)                                                                         |
| -trim     | bool   | Trims transparent edges (default false)                                                                             |
| -tol      | int    | Transparency tolerance for trimming (0-255, default 0)                                                              |
| -same     | bool   | Enable identical image detection (default false)                                                                    |
| -algo     | int    | Packing algorithm (0=Basic, 1=Skyline, 2=MaxRects) (default 1)                                                      |
| -heur     | int    | MaxRects heuristic (0=BestShortSideFit, 1=BestLongSideFit, 2=BestAreaFit, 3=BottomLeft, 4=ContactPoint) (default 0) |

### 🛠️ Unpacking Options

| Parameter | Description                                        |
|-----------|----------------------------------------------------|
| `-u`      | Path to atlas JSON file (required)                 |
| `-img`    | Path to atlas image (optional, inferred from JSON) |
| `-o`      | Output directory (optional, inferred from JSON)    |

### 📦 Examples

```bash
	spritepacker -i <dir> [options]       # Pack mode
	spritepacker -u <json> [options]     # Unpack mode
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
	
	options, _ := (pack.NewOptions().MaxSize(1024, 1024)).Trim(true).Validate()

	atlasInfo, atlasImages, _  := pack.NewPacker(options).PackSprites("./input")
	
	for i := range atlasImages {
		outputPath := filepath.Join("output", atlasInfo.Atlases[i].Name)
		_ = pack.SaveImgAutoExt(outputPath, atlasImages[i], pack.PNG, pack.WithCLV(pack.DefaultCompression))
	}
	
	jsonBytes, _ := json.MarshalIndent(atlasInfo, "", "  ")
	_ = os.WriteFile("output/atlas.json", jsonBytes, os.ModePerm)
}

func Unpack() {
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
