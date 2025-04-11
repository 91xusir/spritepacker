package packer

import (
	"spritepacker/internal/shape"
)

type algoMaxrects struct {
	bin         *shape.Bin // The bin to pack
	allowRotate bool       // Whether to allow rectangle rotation
}

func (algo *algoMaxrects) init(bin *shape.Bin, isRotateEnable bool) {
	algo.bin = bin
	algo.allowRotate = isRotateEnable
}

func (algo *algoMaxrects) packing() *PackResult {

	return &PackResult{}
}
