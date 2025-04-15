/*
Package main provides a command-line tool for sprite atlas packing and unpacking.

This tool supports:
- Creating optimized sprite atlases from individual images
- Unpacking existing sprite atlases back to individual images
- Multiple packing algorithms and heuristics
- Various sprite processing options

Usage:

	spritepacker -i <dir> [options]       # Pack mode
	spritepacker -u <json> [options]      # Unpack mode

Packing Options:

	-I    string  Input directory containing sprite images (required for packing)
	-o    string  Output directory (default "output")
	-maxw int     Maximum atlas width (default 2048)
	-maxh int     Maximum atlas height (default 2048)
	-pad  int     Padding between sprites (default 0)
	-auto         Automatically adjust atlas size (default true)
	-rot          Allow sprite rotation to save space (default false)
	-pot          Force power-of-two atlas dimensions (default false)
	-name string  Base name for output files (default "atlas")
	-sort         Sorts sprites before packing (default true)
	-trim         Trims transparent edges (default false)
	-tol  int     Transparency tolerance for trimming (0-255, default 0)
	-same         Enable identical image detection (default false)
	-algo int     Packing algorithm (0=Basic, 1=Skyline, 2=MaxRects) (default 1)
	-heur int     MaxRects heuristic (0-4, see docs) (default 0)

Unpacking Options:

	-u   string     JSON file to unpack (required for unpacking)
	-img string     Atlas image path (optional, inferred from JSON)
	-o   string     Output directory for unpacked sprites (optional, inferred from JSON)

Examples:

	# Pack sprites with default settings
	spritepacker -i ./sprites -o ./atlases

	# Unpack atlas
	spritepacker -u ./atlases/atlas.json -o ./unpacked

The package includes:
- Multiple packing algorithms (Basic, Skyline, MaxRects)
- Support for common image formats (PNG, JPEG, BMP, TIFF)
- Sprite trimming and optimization features
- JSON metadata generation
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/91xusir/spritepacker/pack"
	"os"
	"path/filepath"
	"runtime/debug"
)

var (
	inputPath      string
	outputPath     string
	unpackJsonPath string
	atlasImgPath   string
	name           string
)

// flagArgs function to parse the command line arguments and populate the options
func flagArgs(opts *pack.Options) error {
	// ---- atlas layout options ----
	maxW := flag.Int("maxw", 2048, "Maximum atlas width (default 2048)")
	maxH := flag.Int("maxh", 2048, "Maximum atlas height (default 2048)")
	autoSize := flag.Bool("auto", true, "Automatically adjust atlas size (default true)")
	padding := flag.Int("pad", 0, "Padding between sprites in pixels (default 0)")
	allowRotate := flag.Bool("rot", false, "Allow sprite rotation to save space (default false)")
	powerOfTwo := flag.Bool("pot", false, "Force atlas size to power of two (default false)")
	name = *flag.String("name", "atlas", "Atlas name (default 'atlas')")

	// ---- sprite processing options ----
	sort := flag.Bool("sort", true, "Sort sprites by Area before packing (default true)")
	trim := flag.Bool("trim", false, "Trim transparent edges from sprites (default false)")
	tolerance := flag.Int("tol", 0, "Tolerance level for trimming (0-255) (default 0)")
	sameDetect := flag.Bool("same", false, "Enable identical image detection (default false)")

	// ---- algorithm settings ----
	algorithm := flag.Int("algo", int(pack.AlgoSkyline), "Packing algorithm: 0=Basic, 1=Skyline, 2=MaxRects (Default: Skyline)")
	heuristic := flag.Int("heur", int(pack.BestShortSideFit), "Heuristic for MaxRects (if used) 0=BestShortSideFit, 1=BestLongSideFit, 2=BestAreaFit, 3=BottomLeftRule, 4=ContactPointRule (Default: BestShortSideFit)")

	// ---- general settings ----
	flag.StringVar(&inputPath, "i", "", "Input directory containing sprite images")
	flag.StringVar(&outputPath, "o", "", "Output directory to save atlases or unpacked sprites")
	flag.StringVar(&unpackJsonPath, "u", "", "Unpack from JSON file")
	flag.StringVar(&atlasImgPath, "img", "", "Atlas image path for unpacking")
	//version
	vFlag := flag.Bool("v", false, "Show version")

	cFlag := flag.Bool("c", false, "compare input and output images")

	flag.Parse()

	if *vFlag {
		fmt.Println("SpritePacker " + pack.Version)
		os.Exit(0)
	}

	if *cFlag {
		diffs := pack.CompareImgFormFolders(inputPath, outputPath)
		if len(diffs) > 0 {
			fmt.Printf("Found %d different images:\n", len(diffs))
		} else {
			fmt.Printf("All images are the same.\n")
		}
		os.Exit(0)
	}

	// apply parsed flags to options
	opts.MaxSize(*maxW, *maxH).
		AutoSize(*autoSize).
		Padding(*padding).
		AllowRotate(*allowRotate).
		PowerOfTwo(*powerOfTwo).
		Sort(*sort).
		Trim(*trim).
		Tolerance(*tolerance).
		SameDetect(*sameDetect).
		Name(name).
		Algorithm(pack.Algorithm(*algorithm)).
		Heuristic(pack.Heuristic(*heuristic))

	return nil
}

func main() {

	info, _ := debug.ReadBuildInfo()
	if info != nil {
		pack.Version = info.Main.Version
	}

	opts := pack.NewOptions()
	check(flagArgs(opts))

	if unpackJsonPath != "" {
		check(pack.UnpackSprites(unpackJsonPath, pack.WithImg(atlasImgPath), pack.WithOutput(outputPath)))
		os.Exit(0)
	}

	if inputPath == "" {
		panic("input path is empty")
	}

	spriteAtlasInfo, atlasImages, err := pack.NewPacker(opts).PackSprites(inputPath)
	check(err)

	for i := range atlasImages {
		filePath := filepath.Join(outputPath, spriteAtlasInfo.Atlases[i].Name)
		check(pack.SaveImgByExt(filePath, atlasImages[i], pack.WithCLV(pack.DefaultCompression)))
	}

	jsonBytes, err := json.MarshalIndent(spriteAtlasInfo, "", "  ")
	check(err)
	check(os.WriteFile(filepath.Join(outputPath, name+".json"), jsonBytes, os.ModePerm))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
