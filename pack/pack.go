package pack

import (
	"fmt"
	"github.com/91xusir/spritepacker/model"
	"github.com/91xusir/spritepacker/utils"
	"image"
	"image/draw"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var (
	Version = "dev"
)

const (
	Repo   = "https://github.com/91xusir/spritepacker"
	Format = "RGBA8888"
)

type Packer struct {
	algo           algo     // interface algo
	option         *Options // Options for packing
	sameDetectInfo utils.SameDetectInfo
	inputDir       string // input path
}

func NewPacker(option *Options) *Packer {
	p := &Packer{option: option}
	switch option.algorithm {
	case AlgoSkyline:
		p.algo = &algoSkyline{}
	case AlgoMaxRects:
		p.algo = &algoMaxrects{}
	default:
		p.algo = &algoBasic{}
	}
	return p
}

// PackRect packs the given rectangles into a bin and returns the result.
func (p *Packer) PackRect(reqRects []model.Rect) []model.Bin {

	var bins []model.Bin
	if len(reqRects) == 0 {
		return bins
	}

	// init algo
	p.algo.init(p.option)

	// sort rects
	if p.option.sort {
		sort.Slice(reqRects, func(i, j int) bool {
			return reqRects[i].Area() > reqRects[j].Area()
		})
	}

	// add padding
	if p.option.padding != 0 {
		for i := range reqRects {
			addPadding(&reqRects[i], p.option.padding)
		}
	}

	bins = p.packInBins(reqRects)

	// remove padding
	if p.option.padding != 0 {
		for i := range bins {
			for j := range bins[i].PackedRects {
				addPadding(&bins[i].PackedRects[j], -p.option.padding)
			}
		}
	}
	return bins
}

func addPadding(r *model.Rect, padding int) {
	r.W += padding
	r.H += padding
}

// PackSprites packs the given sprite images
//
// Parameters:
//   - spritePaths: the paths of the sprite images
//
// Returns:
//   - *AtlasInfo: the sprite atlas info
//   - []image.Image: the atlas images
//   - error
//
// Example:
//
//	spriteAtlas, atlasImages, err := packer.PackSprites("./input")
func (p *Packer) PackSprites(input string) (*model.AtlasInfo, []image.Image, error) {
	spritePaths, err := utils.ListFilePaths(input)
	if err != nil {
		return nil, nil, err
	}

	if p.option.sameDetect {
		spritePaths, p.sameDetectInfo, _ = utils.FindDuplicateFiles(spritePaths)
	}

	// save input dir
	p.inputDir = input

	// create meta
	meta := getMateData()

	// create sprite atlas
	spriteAtlas := &model.AtlasInfo{
		Meta:    meta,
		Atlases: make([]model.Atlas, 0),
	}
	// get image rects and src rects and trimmed rects
	reqRects, srcRects, trimmedRectMap := p.getImageRects(spritePaths)

	// pack rects
	bins := p.PackRect(reqRects)

	// generate atlases info
	//AtlasInfo->
	//			Meta
	//			[]Atlases->
	//					Name
	//					Size
	//					[]Sprites->
	//							FileName
	//							Frame
	//							SrcRect
	//							TrimmedRect
	//							Rotated
	//							Trimmed
	for i, bin := range bins {

		atlasSize := model.Size{W: bin.W, H: bin.H}

		// if power of two
		if p.option.powerOfTwo {
			atlasSize = atlasSize.PowerOfTwo()
		}

		// create atlas
		var atlasName string
		if len(bins) == 1 {
			atlasName = p.option.name + p.option.imgExt
		} else {
			atlasName = fmt.Sprintf("%s_%d%s", p.option.name, i, p.option.imgExt)
		}
		atlas := model.Atlas{
			Name:    atlasName,
			Size:    atlasSize,
			Sprites: make([]model.Sprite, 0, len(bin.PackedRects)),
		}

		//fmt.Printf("len %d \n", len(bin.PackedRects))
		//for _, rect := range bin.PackedRects {
		//	fmt.Printf("rect %v \n", rect)
		//}

		for _, rect := range bin.PackedRects {
			// create sprite
			baseName := filepath.Base(spritePaths[rect.Id])
			sprite := model.Sprite{
				FileName:    baseName,
				Frame:       rect,
				SrcRect:     srcRects[rect.Id],
				TrimmedRect: trimmedRectMap[rect.Id],
				Rotated:     rect.IsRotated,
				Trimmed:     p.option.trim,
			}
			atlas.Sprites = append(atlas.Sprites, sprite)

			// if same detect
			// try to find the same file in the same directory
			if p.option.sameDetect {
				if dupPaths, ok := p.sameDetectInfo.BaseToDupesName[baseName]; ok {
					//fmt.Printf("Found duplicate files: %v\n", dupPaths)
					for _, dupPath := range dupPaths {
						s := sprite.Clone()
						s.FileName = dupPath
						atlas.Sprites = append(atlas.Sprites, s)
					}
				}
			}
		}

		spriteAtlas.Atlases = append(spriteAtlas.Atlases, atlas)
	}
	images, err := p.createAtlasImages(spriteAtlas)
	if err != nil {
		return spriteAtlas, nil, err
	}
	return spriteAtlas, images, nil
}

func (p *Packer) packInBins(reqRects []model.Rect) []model.Bin {
	var bins []model.Bin
	remainingRects := reqRects
	// loop until all rects are packed
	for len(remainingRects) > 0 {
		// reset algo
		p.algo.reset(p.option.maxW, p.option.maxH)

		// Try packing the remaining rectangles into a new bin
		packedRects, unpackedRects := p.algo.packing(remainingRects)

		if len(packedRects) == 0 {
			//If no rectangle can be packed, it may be an algorithm problem or the rectangle is too large
			//Log and jump out of the loop to avoid infinite loops
			_, _ = fmt.Fprintf(os.Stderr, "Warning: Unable to pack the remaining %d rectangles", len(remainingRects))
			break
		}

		// calculates the total area of the packed rectangle
		totalArea := 0
		for _, rect := range packedRects {
			totalArea += rect.W * rect.H
		}

		// If there are no unpacked rectangles and autosize is enabled, try optimizing the bin size
		if len(unpackedRects) == 0 && p.option.autoSize {
			// calculates the minimum side length of the square
			minSide := int(math.Ceil(math.Sqrt(float64(totalArea))))

			// set the scope of your search
			low := minSide
			high := utils.MaxInt(p.option.maxH, p.option.maxW)

			var bestSize int
			var bestResult []model.Rect
			found := false

			// Try to find the smallest feasible square size by binocular
			for low <= high {
				mid := (low + high) / 2
				p.algo.reset(mid, mid)
				packs, unpacks := p.algo.packing(packedRects)

				if len(unpacks) == 0 {
					bestResult = packs
					bestSize = mid
					found = true
					high = mid - 1
				} else {
					low = mid + 1
				}
			}

			if found {
				// create a bin using the optimal size found
				bin := model.NewBin(bestSize, bestSize, bestResult)
				bin.UsedArea = totalArea
				bins = append(bins, bin)
			} else {
				// If you can't find the optimal size, use the original size
				bin := model.NewBin(p.option.maxW, p.option.maxH, packedRects)
				bin.UsedArea = totalArea
				bins = append(bins, bin)
			}
		} else {
			// If there are unpacked rectangles or autosize not enabled, use the original size
			bin := model.NewBin(p.option.maxW, p.option.maxH, packedRects)
			bin.UsedArea = totalArea
			bins = append(bins, bin)
		}

		// Update the remaining rectangles that need to be packaged
		remainingRects = unpackedRects
	}

	return bins
}

func (p *Packer) getImageRects(filePaths []string) ([]model.Rect, []model.Size, map[int]model.Rect) {
	reqRects := make([]model.Rect, 0)
	srcRects := make([]model.Size, len(filePaths))
	trimmedRectMap := make(map[int]model.Rect)
	for i, fileName := range filePaths {
		file, err := os.Open(fileName)
		if err != nil {
			continue // Skip unreadable files
		}
		//if i == 1 || i == 2 {
		//	// continue may cause an empty rect to be passed in, resulting in an extra rect with ID default of 0
		//	// so use reqRects.append replace reqRects[i]
		//	// fix on 2025/4.14
		//	continue
		//}
		if p.option.trim {
			src, _, err := image.Decode(file)
			file.Close()
			if err != nil {
				continue // Skip non-image files
			}
			srcSize := model.Size{
				W: src.Bounds().Dx(),
				H: src.Bounds().Dy(),
			}
			trimRect := utils.GetOpaqueBounds(src, p.option.tolerance)
			trimmedRect := model.NewRectByPosAndSize(
				trimRect.Min.X,
				trimRect.Min.Y,
				trimRect.Dx(),
				trimRect.Dy(),
			)
			srcRects[i] = srcSize
			reqRects = append(reqRects, model.NewRectBySizeAndId(trimRect.Dx(), trimRect.Dy(), i))
			trimmedRectMap[i] = trimmedRect
		} else {
			cfg, _, err := image.DecodeConfig(file)
			file.Close()
			if err != nil {
				continue // Skip non-image files
			}
			srcSize := model.Size{
				W: cfg.Width,
				H: cfg.Height,
			}
			srcRects[i] = srcSize
			reqRects = append(reqRects, model.NewRectBySizeAndId(cfg.Width, cfg.Height, i))
		}
	}

	return reqRects, srcRects, trimmedRectMap
}

func (p *Packer) createAtlasImages(atlas *model.AtlasInfo) ([]image.Image, error) {
	var atlasImages = make([]image.Image, len(atlas.Atlases))
	for i := range atlas.Atlases {
		atlasSize := atlas.Atlases[i].Size
		// create atlas image
		atlasImg := image.NewNRGBA(image.Rect(0, 0, atlasSize.W, atlasSize.H))
		for j := range atlas.Atlases[i].Sprites {
			sprite := atlas.Atlases[i].Sprites[j]
			trimmedRect := sprite.TrimmedRect
			srcLeftTopPoint := image.Point{
				X: trimmedRect.X,
				Y: trimmedRect.Y,
			}
			// if same detect
			if p.option.sameDetect {

				if _, ok := p.sameDetectInfo.DupeToBaseName[sprite.FileName]; ok {
					//fmt.Printf("same detect %s \n", sprite.FileName)
					continue
				}
			}

			spriteImg, err := utils.LoadImg(filepath.Join(p.inputDir, sprite.FileName))
			if err != nil {
				return nil, err
			}
			// if rotated
			if sprite.Rotated {
				spriteImg = utils.Rotate270(spriteImg)
				srcH := sprite.SrcRect.H
				newX := srcH - trimmedRect.Y - trimmedRect.H
				newY := trimmedRect.X
				srcLeftTopPoint.X = newX
				srcLeftTopPoint.Y = newY
			}
			ditPosition := sprite.Frame.ToImageRect()
			draw.Draw(atlasImg, ditPosition, spriteImg, srcLeftTopPoint, draw.Src)
			atlasImages[i] = atlasImg
		}
	}
	return atlasImages, nil
}

func getMateData() model.Meta {
	return model.Meta{
		Repo:      Repo,
		Format:    Format,
		Version:   Version,
		Timestamp: time.Now().Format(time.DateTime),
	}
}
