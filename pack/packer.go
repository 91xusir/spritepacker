package pack

import (
	"fmt"
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

type packer struct {
	algo   algo     // interface algo
	option *Options // Options for packing
}

func NewPacker(option *Options) *packer {
	p := &packer{option: option}
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

//func (p *packer) PackImage() *PackedResult {
//	return p.algo.packing()
//}

func (p *packer) PackRect(reqRects []Rect) PackedResult {
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

	p.algo.init(p.option)

	return p.algo.packing(reqRects)
}
