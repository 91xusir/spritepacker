package packer

import (
	"fmt"
	"spritepacker/internal/shape"
)

// PackResult represents the result of packing rectangles into a bin.
type PackResult struct {
	PackedRects   []*shape.RectPacked // Packed rectangle list
	UnpackedRects []*shape.Rect       // Unpacked rectangle list
	TotalArea     int                 // Total area
	FillRate      float64             // Fill rate
}
func (p *PackResult) String() string {
	return fmt.Sprintf("PackedRects: %v \n UnpackedRects: %v \n TotalArea: %d, FillRate: %.2f%%", p.PackedRects, p.UnpackedRects, p.TotalArea, p.FillRate)
}

type packer struct {
	algo   algorithm // interface algorithm
	option *options  // options for packing
}

func NewPacker(bin *shape.Bin, option *options) *packer {
	p := &packer{option: option}
	switch option.Algorithm {
	case AlgoBasic:
		p.algo = &basicAlgo{}
	case AlgoSkyline:
		p.algo = &algoSkyline{}
	}
	p.algo.init(bin, option.AllowRotate)
	return p
}

func (p *packer) Packing() *PackResult {
	// TODO: add sort
	// TODO: add trim
	return p.algo.packing()
}
