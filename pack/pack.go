package pack

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"sort"
	"time"
)

const (
	Version = "0.0.1"
	Repo    = "https://github.com/91xusir/spritepacker"
	Format  = "RGBA8888"
)

type Packer struct {
	algo   algo     // interface algo
	option *Options // Options for packing
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
func (p *Packer) PackRect(reqRects []Rect) []Bin {

	var bins []Bin
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
				removePadding(&bins[i].PackedRects[j], p.option.padding)
			}
		}
	}
	return bins
}

// PackSprites packs the given sprite images
//
// Parameters:
//   - spritePaths: the paths of the sprite images
//
// Returns:
//   - *SpriteAtlas: the sprite atlas info
//   - []image.Image: the atlas images
//   - error
func (p *Packer) PackSprites(spritePaths []string) (*SpriteAtlas, []image.Image, error) {
	// create meta
	meta := Meta{
		Repo:      Repo,
		Format:    Format,
		Version:   Version,
		Timestamp: time.Now().Format(time.DateTime),
	}
	// create sprite atlas
	spriteAtlas := &SpriteAtlas{
		Meta:    meta,
		Atlases: make([]Atlas, 0),
	}
	// get image rects and src rects and trimmed rects
	reqRects, srcRects, trimmedRectMap := p.getImageRects(spritePaths)
	bins := p.PackRect(reqRects)
	for i, bin := range bins {
		atlasSize := Size{W: bin.W, H: bin.H}
		// if power of two
		if p.option.powerOfTwo {
			atlasSize = atlasSize.PowerOfTwo()
		}
		// create atlas
		atlas := Atlas{
			Name:    fmt.Sprintf("%s_%d", p.option.name, i),
			Size:    atlasSize,
			Sprites: make([]Sprite, len(bin.PackedRects)),
		}
		for j, rect := range bin.PackedRects {
			//fmt.Printf("%+v, srcPath: %s\n", rect.Id, spritePaths[rect.Id])
			atlas.Sprites[j] = Sprite{
				Filepath:    spritePaths[rect.Id], // always use filepath, not filename
				Frame:       rect,
				SrcRect:     srcRects[rect.Id],
				TrimmedRect: trimmedRectMap[rect.Id],
				Rotated:     rect.IsRotated,
				Trimmed:     p.option.trim,
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

func (p *Packer) packInBins(reqRects []Rect) []Bin {
	var bins []Bin
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
			high := MaxInt(p.option.maxH, p.option.maxW)

			var bestSize int
			var bestResult []Rect
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
				bin := NewBin(bestSize, bestSize, bestResult)
				bin.UsedArea = totalArea
				bins = append(bins, bin)
			} else {
				// If you can't find the optimal size, use the original size
				bin := NewBin(p.option.maxW, p.option.maxH, packedRects)
				bin.UsedArea = totalArea
				bins = append(bins, bin)
			}
		} else {
			// If there are unpacked rectangles or autosize not enabled, use the original size
			bin := NewBin(p.option.maxW, p.option.maxH, packedRects)
			bin.UsedArea = totalArea
			bins = append(bins, bin)
		}

		// Update the remaining rectangles that need to be packaged
		remainingRects = unpackedRects
	}

	return bins
}

func (p *Packer) getImageRects(fileNames []string) ([]Rect, []Size, map[int]Rect) {
	reqRects := make([]Rect, len(fileNames))
	srcRects := make([]Size, len(fileNames))
	trimmedRectMap := make(map[int]Rect)
	for i, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println(err)
			continue // Skip unreadable files
		}
		if p.option.trim {
			src, _, err := image.Decode(file)
			file.Close()
			if err != nil {
				fmt.Println(err)
				continue // Skip non-image files
			}
			srcSize := Size{
				W: src.Bounds().Dx(),
				H: src.Bounds().Dy(),
			}
			trimRect := GetOpaqueBounds(src, p.option.tolerance)
			trimmedRect := NewRectByPosAndSize(
				trimRect.Min.X,
				trimRect.Min.Y,
				trimRect.Dx(),
				trimRect.Dy(),
			)
			srcRects[i] = srcSize
			reqRects[i] = NewRectBySizeAndId(trimRect.Dx(), trimRect.Dy(), i)
			trimmedRectMap[i] = trimmedRect
		} else {
			cfg, _, err := image.DecodeConfig(file)
			file.Close()
			if err != nil {
				continue // Skip non-image files
			}
			srcSize := Size{
				W: cfg.Width,
				H: cfg.Height,
			}
			srcRects[i] = srcSize
			reqRects[i] = NewRectBySizeAndId(cfg.Width, cfg.Height, i)
		}
	}

	return reqRects, srcRects, trimmedRectMap
}

func (p *Packer) createAtlasImages(atlas *SpriteAtlas) ([]image.Image, error) {
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
			//read sprite image
			spriteImg, err := LoadImg(sprite.Filepath)
			if err != nil {
				return nil, err
			}
			// if rotated
			if sprite.Rotated {
				spriteImg = Rotate270(spriteImg)
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
