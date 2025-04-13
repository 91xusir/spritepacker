package pack

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
func (p *Packer) PackRect(reqRects []Rect) PackedResult {
	if len(reqRects) == 0 {
		return PackedResult{
			Bin:           NewBin(p.option.maxW, p.option.maxH, make([]PackedRect, 0), 0, 0),
			UnpackedRects: make([]Rect, 0),
		}
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

	// get result
	result := p.algo.packing(reqRects)

	// if autosize
	if p.option.autoSize {
		var ok bool
		result, ok = p.autosize(result)
		if !ok {
			fmt.Println("warning: cannot autosize bin size")
		}
	}

	// remove padding
	if p.option.padding != 0 {
		for i := range result.Bin.PackedRects {
			removePadding(&result.Bin.PackedRects[i].Rect, p.option.padding)
		}
	}

	return result
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

	remainingRects := reqRects

	atlasIndex := 0
	for len(remainingRects) > 0 {
		// run packing
		packedResult := p.PackRect(remainingRects)
		if len(packedResult.Bin.PackedRects) == 0 {
			break
		}
		atlasSize := Size{W: packedResult.Bin.W, H: packedResult.Bin.H}
		// if power of two
		if p.option.powerOfTwo {
			atlasSize = atlasSize.PowerOfTwo()
		}
		// create atlas
		atlas := Atlas{
			Name:    fmt.Sprintf("%s_%d", p.option.name, atlasIndex),
			Size:    atlasSize,
			Sprites: make([]Sprite, len(packedResult.Bin.PackedRects)),
		}
		for i, rect := range packedResult.Bin.PackedRects {
			atlas.Sprites[i] = Sprite{
				Filepath:    spritePaths[rect.Id], // always use filepath, not filename
				Frame:       rect.ToRectangle(),
				SrcRect:     srcRects[rect.Id],
				TrimmedRect: trimmedRectMap[rect.Id],
				Rotated:     rect.IsRotated,
				Trimmed:     p.option.trim,
			}
		}
		spriteAtlas.Atlases = append(spriteAtlas.Atlases, atlas)
		remainingRects = packedResult.UnpackedRects
		atlasIndex++
	}
	images, err := p.createAtlasImages(spriteAtlas)
	if err != nil {
		return spriteAtlas, nil, err
	}
	return spriteAtlas, images, nil
}

func (p *Packer) autosize(result PackedResult) (PackedResult, bool) {
	if len(result.UnpackedRects) != 0 {
		// unpacked packedRects not empty means no space left
		return result, false
	}
	totalArea := result.Bin.UsedArea
	if totalArea == 0 {
		return result, false
	}
	packedRects := result.Bin.PackedRects
	if len(packedRects) == 0 {
		return result, false
	}
	// calculate min side length of square
	minSide := int(math.Ceil(math.Sqrt(float64(totalArea))))
	// set search range
	low := minSide
	high := MaxInt(result.Bin.W, result.Bin.H) // set high to 2 * minSide to ensure we can find a solution
	// copy packedRects to reqRects
	reqRects := make([]Rect, len(packedRects))
	for i, rect := range packedRects {
		reqRects[i] = NewRectById(rect.W, rect.H, rect.Id)
	}
	var bestResult PackedResult
	found := false
	// try binary search
	for low <= high {
		mid := (low + high) / 2
		p.algo.reset(mid, mid)
		re := p.algo.packing(reqRects)
		if len(re.UnpackedRects) == 0 {
			bestResult = re
			found = true
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	if !found {
		return result, false
	}
	return bestResult, true
}
func (p *Packer) getImageRects(fileNames []string) ([]Rect, []Size, map[int]Rectangle) {
	reqRects := make([]Rect, 0, len(fileNames))
	srcRects := make([]Size, 0, len(fileNames))
	trimmedRectMap := make(map[int]Rectangle)

	for _, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			continue // Skip unreadable files
		}
		if p.option.trim {
			src, _, err := image.Decode(file)
			file.Close()
			if err != nil {
				continue // Skip non-image files
			}
			srcSize := Size{
				W: src.Bounds().Dx(),
				H: src.Bounds().Dy(),
			}
			trimRect := GetOpaqueBounds(src, p.option.tolerance)
			trimmedRect := Rectangle{
				X: trimRect.Min.X,
				Y: trimRect.Min.Y,
				W: trimRect.Dx(),
				H: trimRect.Dy(),
			}
			srcRects = append(srcRects, srcSize)
			reqRects = append(reqRects, NewRectById(trimRect.Dx(), trimRect.Dy(), len(reqRects)))
			trimmedRectMap[len(reqRects)-1] = trimmedRect
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
			srcRects = append(srcRects, srcSize)
			reqRects = append(reqRects, NewRectById(cfg.Width, cfg.Height, len(reqRects)))
		}
	}

	return reqRects, srcRects, trimmedRectMap
}

func (p *Packer) createAtlasImages(atlas *SpriteAtlas) ([]image.Image, error) {
	var atlasImages []image.Image = make([]image.Image, len(atlas.Atlases))
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

type unpackedOpts struct {
	atlasImgPath string
	outputPath   string
}
type UnpackOpts func(*unpackedOpts)

func WithAtlasImgPath(atlasImgPath string) UnpackOpts {
	if atlasImgPath == "" {
		return func(opts *unpackedOpts) {
			return
		}
	}
	return func(opts *unpackedOpts) {
		opts.atlasImgPath = atlasImgPath
	}
}
func WithOutputPath(outputPath string) UnpackOpts {
	if outputPath == "" {
		return func(opts *unpackedOpts) {
			return
		}
	}
	return func(opts *unpackedOpts) {
		opts.outputPath = outputPath
	}
}

func UnpackSprites(jsonPath string, fn ...UnpackOpts) error {
	opts := &unpackedOpts{
		atlasImgPath: filepath.Dir(jsonPath),
		outputPath:   filepath.Dir(jsonPath),
	}
	for _, f := range fn {
		f(opts)
	}

	// make sure outputPath exists
	if err := os.MkdirAll(opts.outputPath, os.ModePerm); err != nil {
		return err
	}

	jsonData, err := os.ReadFile(jsonPath)
	var atlasInfo SpriteAtlas
	err = json.Unmarshal(jsonData, &atlasInfo)
	if err != nil {
		return err
	}

	baseNames := make([]string, len(atlasInfo.Atlases))
	for i := range atlasInfo.Atlases {
		baseNames[i] = strings.TrimSuffix(filepath.Base(atlasInfo.Atlases[i].Name), filepath.Ext(atlasInfo.Atlases[i].Name))
	}

	extStr := []string{".png", ".jpg", ".jpeg", ".bmp", ".tiff"}

	for i, baseName := range baseNames {
		var imgFilePath string
		found := false
		for _, ext := range extStr {
			tryPath := filepath.Join(opts.atlasImgPath, baseName+ext)
			if _, err := os.Stat(tryPath); err == nil {
				imgFilePath = tryPath
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("image file for atlas %s not found", baseName)
		}
		atlasImg, err := LoadImg(imgFilePath)
		if err != nil {
			return fmt.Errorf("failed to load image %s: %v", imgFilePath, err)
		}
		for j := range atlasInfo.Atlases[i].Sprites {
			sprite := atlasInfo.Atlases[i].Sprites[j]
			outputPath := filepath.Join(opts.outputPath, filepath.Base(sprite.Filepath))
			subImg := image.NewNRGBA(image.Rect(0, 0, sprite.Frame.W, sprite.Frame.H))
			srcLeftTopPoint := image.Point{
				X: sprite.Frame.X,
				Y: sprite.Frame.Y,
			}
			draw.Draw(subImg, subImg.Bounds(), atlasImg, srcLeftTopPoint, draw.Src)
			// if rotated
			if sprite.Rotated {
				subImg = Rotate90(subImg)
			}
			// if trimmed
			if sprite.Trimmed {
				img := image.NewNRGBA(image.Rect(0, 0, sprite.SrcRect.W, sprite.SrcRect.H))
				destRect := image.Rect(
					sprite.TrimmedRect.X,
					sprite.TrimmedRect.Y,
					sprite.TrimmedRect.X+subImg.Bounds().Dx(),
					sprite.TrimmedRect.Y+subImg.Bounds().Dy(),
				)
				draw.Draw(img, destRect, subImg, image.Point{}, draw.Src)
				subImg = img
			}
			err := SaveImgExt(outputPath, subImg)
			if err != nil {
				return fmt.Errorf("failed to save image %s: %v", outputPath, err)
			}
		}
	}
	return nil
}
