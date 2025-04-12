package pack

import (
	"fmt"
	"image"
	"log"
	"os"
	"sort"
)

// PackedResult represents the result of packing rectangles into a bin.
type PackedResult struct {
	Bin           Bin
	UnpackedRects []Rect
}

func (r PackedResult) String() string {
	return fmt.Sprintf("PackedResult{Bin:%s, UnpackedRects:%v}", r.Bin, r.UnpackedRects)
}

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

func (p *Packer) PackRect(reqRects []Rect) PackedResult {
	if len(reqRects) == 0 {
		return PackedResult{
			Bin:           NewBin(p.option.maxW, p.option.maxH, make([]PackedRect, 0), 0, 0),
			UnpackedRects: make([]Rect, 0),
		}
	}
	if p.option.sort {
		sort.Slice(reqRects, func(i, j int) bool {
			return reqRects[i].Area() > reqRects[j].Area()
		})
	}
	// if autosize
	if p.option.autoSize {
		//
	}

	p.algo.init(p.option)

	return p.algo.packing(reqRects)
}

func (p *Packer) PackSprites() error {
	atlasesName :=GetLastFolderName(p.option.inputDir)
	spritePaths, err := GetFilesInDirectory(p.option.inputDir)
	if err != nil {
		return err
	}
	reqRects, srcRects, trimmedRectMap := GetImageRects(spritePaths, p.option.trim, p.option.tolerance)

	atlases := make([]Atlas, 0)

	packedResult := p.PackRect(reqRects)
	unPackedRects := packedResult.UnpackedRects
	if len(unPackedRects) == 0 {
		atlases = append(atlases, Atlas{
			Name: atlasesName,
			Size: Size{W: p.algo, H: p.option.maxH},
			Sprites: ,
		})
	}
	atlases = append(atlases, Atlas{})
	for len(unPackedRects) > 0 {
		packedResult = p.PackRect(unPackedRects)
		unPackedRects = packedResult.UnpackedRects
	}

	return nil
}

// GetImageRects reads images from directory and returns their dimensions
func GetImageRects(fileNames []string, trim bool, tolerance uint8) ([]Rect, []Size, map[int]Rectangle) {
	reqRects := make([]Rect, len(fileNames))
	srcRects := make([]Size, len(fileNames))
	trimmedRectMap := make(map[int]Rectangle)
	for id, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			continue // Skip unreadable files
		}
		if trim {
			src, _, err := image.Decode(file)
			file.Close()
			if err != nil {
				continue // Skip non-image files
			}
			srcRects[id] = Size{
				W: src.Bounds().Dx(),
				H: src.Bounds().Dy(),
			}
			trimRect := GetOpaqueBounds(src, tolerance)
			trimmedRectMap[id] = Rectangle{
				X: trimRect.Min.X,
				Y: trimRect.Min.Y,
				W: trimRect.Dx(),
				H: trimRect.Dy(),
			}
			reqRects[id] = NewRectById(trimRect.Dx(), trimRect.Dy(), id)

		} else {
			cfg, _, err := image.DecodeConfig(file)
			file.Close()
			if err != nil {
				continue // Skip non-image files
			}
			srcRects[id] = Size{
				W: cfg.Width,
				H: cfg.Height,
			}
			reqRects[id] = NewRectById(cfg.Width, cfg.Height, id)
		}
	}

	return reqRects, srcRects, trimmedRectMap
}

//
//func processImages(paths []string, options *Options) ([]Rect, []image.Rectangle, error) {
//
//	sourceRects := make([]image.Rectangle, len(paths))
//	sizes := make([]Rect, len(paths))
//	// 创建错误通道
//	errChan := make(chan error, len(paths))
//	// 创建互斥锁保护对errChan的并发访问
//	var mu sync.Mutex
//	// 使用Parallel函数并行处理图片
//	parallel(0, len(paths), func(i int) {
//		path := paths[i]
//		file, err := os.Open(path)
//		if err != nil {
//			mu.Lock()
//			errChan <- err
//			mu.Unlock()
//			return
//		}
//		if options.trim {
//			// 完全解码图片以分析透明区域
//			src, err := imaging.Decode(file)
//			file.Close()
//			if err != nil {
//				mu.Lock()
//				errChan <- fmt.Errorf("无法解码图片 %s: %v", path, err)
//				mu.Unlock()
//				return
//			}
//			// 获取原始尺寸
//			origBounds := src.Bounds()
//			sourceRects[i] = origBounds
//			// 获取透明边界区域
//			trimRect := GetOpaqueBounds(src, options.tolerance)
//			sizes[i] = rectpack.NewSize2DByID(i, trimRect.Dx(), trimRect.Dy())
//			sourceRects[i] = trimRect
//		} else {
//			// 只解码图片头部以获取尺寸信息
//			cfg, _, err := image.DecodeConfig(file)
//			file.Close()
//			if err != nil {
//				mu.Lock()
//				errChan <- fmt.Errorf("无法解码图片 %s: %v", path, err)
//				mu.Unlock()
//				return
//			}
//			// 创建尺寸对象，使用索引作为ID
//			sizes[i] = NewSize2DByID(i, cfg.Width, cfg.Height)
//		}
//	})
//
//	// 检查是否有错误
//	close(errChan)
//	for err := range errChan {
//		if err != nil {
//			return nil, nil, err
//		}
//	}
//
//	return sizes, sourceRects, nil
//}
