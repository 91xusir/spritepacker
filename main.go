/*
Package main provides a command-line tool for sprite atlas packing and unpacking.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/91xusir/spritepacker/export"
	"github.com/91xusir/spritepacker/pack"
	"github.com/91xusir/spritepacker/utils"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

var (
	inputPath      string
	outputPath     string
	unpackJsonPath string
	atlasImgPath   string
	name           string
	infoFormat     string
	imgFormat      string
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
	// ---- sprite processing options ----
	sort := flag.Bool("sort", true, "Sort sprites by Area before packing (default true)")
	trim := flag.Bool("trim", false, "Trim transparent edges from sprites (default false)")
	tolerance := flag.Int("tol", 0, "Tolerance level for trimming (0-255) (default 0)")
	sameDetect := flag.Bool("same", false, "Enable identical image detection (default false)")
	// ---- algorithm settings ----
	algorithm := flag.Int("algo", int(pack.AlgoSkyline), "Packing algorithm: 0=Basic, 1=Skyline, 2=MaxRects (Default: Skyline)")
	heuristic := flag.Int("heur", int(pack.BestShortSideFit), "Heuristic for MaxRects (if used) 0=BestShortSideFit, 1=BestLongSideFit, 2=BestAreaFit, 3=BottomLeftRule, 4=ContactPointRule (Default: BestShortSideFit)")
	// ---- general settings ----
	flag.StringVar(&name, "name", "atlas", "Atlas name (default 'atlas')")
	flag.StringVar(&inputPath, "i", "", "Input directory containing sprite images")
	flag.StringVar(&outputPath, "o", "", "Output directory to save atlases or unpacked sprites")
	flag.StringVar(&unpackJsonPath, "u", "", "Unpack from JSON file")
	flag.StringVar(&atlasImgPath, "img", "", "Atlas image path for unpacking")
	flag.StringVar(&infoFormat, "f1", "json", "Atlas info format  (default 'json')")
	flag.StringVar(&imgFormat, "f2", "png", "Atlas image format (default 'png')")
	//version
	vFlag := flag.Bool("v", false, "Show version")

	cFlag := flag.Bool("c", false, "compare input and output images")

	flag.Parse()

	if *vFlag {
		fmt.Println("SpritePacker " + pack.Version)
		os.Exit(0)
	}

	if *cFlag {
		diffs := utils.CompareImgFormFolders(inputPath, outputPath)
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
		ImgExt(imgFormat).
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
		check(pack.UnpackSprites(unpackJsonPath, pack.WithImgInput(atlasImgPath), pack.WithOutput(outputPath)))
		os.Exit(0)
	}

	args := flag.Args()
	// if no input path is specified, use the first argument as input path
	if len(args) > 0 && inputPath == "" {
		inputPath = args[0]
		f, err := os.Stat(inputPath)
		check(err)
		// if input path is a file, unpack it
		if !f.IsDir() {
			check(pack.UnpackSprites(inputPath, pack.WithImgInput(atlasImgPath), pack.WithOutput(outputPath)))
			os.Exit(0)
		}
		// use default options if output path is not specified
		opts.Default()
		name = utils.GetLastFolderName(inputPath)
		opts.Name(name)
	}

	if inputPath == "" {
		panic("input path is empty")
	}
	fmt.Printf("input path: %s\n", inputPath)
	fmt.Printf("output path: %s\n", outputPath)
	fmt.Printf("info format: %s\n", infoFormat)
	fmt.Printf("image format: %s\n", imgFormat)

	spriteAtlasInfo, atlasImages, err := pack.NewPacker(opts).PackSprites(inputPath)
	check(err)

	for i := range atlasImages {
		filePath := filepath.Join(outputPath, spriteAtlasInfo.Atlases[i].Name)
		check(utils.SaveImgByExt(filePath, atlasImages[i], utils.WithCLV(utils.DefaultCompression)))
	}
	exporter := export.NewExportManager().Init()
	_ = exporter.Export(filepath.Join(outputPath, name+dotFormat(infoFormat)), spriteAtlasInfo)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func dotFormat(format string) string {
	return "." + strings.TrimPrefix(format, ".")
}
