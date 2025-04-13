/*
Package main provides a command-line tool for sprite atlas packing and unpacking.

This tool supports:
- Creating optimized sprite atlases from individual images
- Unpacking existing sprite atlases back to individual images
- Multiple packing algorithms and heuristics
- Various sprite processing options

Usage:

	spritepacker -input <dir> [options]       # Pack mode
	spritepacker -unpack <json> [options]     # Unpack mode

Packing Options:

	-input string      Input directory containing sprite images (required for packing)
	-output string     Output directory (default "output")
	-maxw int          Maximum atlas width (default 2048)
	-maxh int          Maximum atlas height (default 2048)
	-padding int       Padding between sprites (default 2)
	-rotate            Allow sprite rotation to save space
	-pot               Force power-of-two atlas dimensions
	-pma               Use premultiplied alpha
	-name string       Base name for output files (default "atlas")
	-sort              Sorts sprites before packing (default true)
	-trim              Trims transparent edges (default true)
	-tolerance int     Transparency tolerance for trimming (0-255, default 1)
	-same              Enable identical image detection (default true)
	-algo int          Packing algorithm (0=Basic, 1=Skyline, 2=MaxRects)
	-heuristic int     MaxRects heuristic (0-4, see docs)

Unpacking Options:

	-unpack string     JSON file to unpack (required for unpacking)
	-img string        Atlas image path (optional, will search by name)
	-output string     Output directory for unpacked sprites

Examples:

	# Pack sprites with default settings
	spritepacker -input ./sprites -output ./atlases

	# Pack with custom settings
	spritepacker -input ./sprites -maxw 1024 -maxh 1024 -padding 4 -rotate

	# Unpack atlas
	spritepacker -unpack ./atlases/atlas.json -output ./unpacked

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
	"github.com/91xusir/spritepacker/pack"
	"os"
	"path/filepath"
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
	maxW := flag.Int("maxw", 2048, "Maximum atlas width")
	maxH := flag.Int("maxh", 2048, "Maximum atlas height")
	autoSize := flag.Bool("autosize", false, "Automatically adjust atlas size based on input")
	padding := flag.Int("padding", 2, "Padding between sprites in pixels")
	allowRotate := flag.Bool("rotate", false, "Allow sprite rotation to save space")
	powerOfTwo := flag.Bool("pot", false, "Force atlas size to power of two")
	preMultipliedAlpha := flag.Bool("pma", false, "Use premultiplied alpha for output")
	name = *flag.String("name", "atlas", "Atlas name")

	// ---- sprite processing options ----
	sort := flag.Bool("sort", true, "Sort sprites by Area before packing")
	trim := flag.Bool("trim", true, "Trim transparent edges from sprites")
	tolerance := flag.Int("tolerance", 1, "Tolerance level for trimming (0-255)")
	sameDetect := flag.Bool("same", true, "Enable identical image detection")

	// ---- algorithm settings ----
	algorithm := flag.Int("algo", int(pack.AlgoSkyline), "Packing algorithm: 0=Basic, 1=Skyline, 2=MaxRects")
	heuristic := flag.Int("heuristic", int(pack.BestShortSideFit), "Heuristic for MaxRects (if used) 0=BestShortSideFit, 1=BestLongSideFit, 2=BestAreaFit, 3=BottomLeftRule, 4=ContactPointRule")

	// ---- general settings ----
	flag.StringVar(&inputPath, "input", "", "Input directory containing sprite images")
	flag.StringVar(&outputPath, "output", "output", "Output directory to save atlases or unpacked sprites")
	flag.StringVar(&unpackJsonPath, "unpack", "", "Unpack from JSON file")
	flag.StringVar(&atlasImgPath, "img", "", "Atlas image path for unpacking")

	flag.Parse()

	// apply parsed flags to options
	opts.MaxSize(*maxW, *maxH).
		AutoSize(*autoSize).
		Padding(*padding).
		AllowRotate(*allowRotate).
		PowerOfTwo(*powerOfTwo).
		PmAlpha(*preMultipliedAlpha).
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
	opts := pack.NewOptions()
	check(flagArgs(opts))
	if unpackJsonPath != "" {
		check(pack.UnpackSprites(unpackJsonPath, pack.WithAtlasImgPath(atlasImgPath), pack.WithOutputPath(outputPath)))
		os.Exit(0)
	}
	if inputPath == "" {
		panic("input path is empty")
	}
	spritePaths, err := pack.GetFilesInDirectory(inputPath)
	check(err)
	spriteAtlasInfo, atlasImages, err := pack.NewPacker(opts).PackSprites(spritePaths)
	check(err)
	for i := range atlasImages {
		outputPath := filepath.Join(outputPath, spriteAtlasInfo.Atlases[i].Name)
		check(pack.SaveImg(outputPath, atlasImages[i], pack.PNG, pack.WithCLV(pack.DefaultCompression)))
	}
	jsonBytes, err := json.MarshalIndent(spriteAtlasInfo, "", "  ")
	check(err)
	check(os.WriteFile(filepath.Join(outputPath, name), jsonBytes, os.ModePerm))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
