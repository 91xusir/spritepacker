package packer

import (
	"spritepacker/internal/shape"
)
const (
	AlgoBasic = iota
	AlgoSkyline
	AlgoMaxRectsBaf
	AlgoMaxRectsBlsf
	AlgoMaxRectsBssf
)

// algorithm is the interface that wraps the Pack method.
type algorithm interface {
	init(bin *shape.Bin, isRotateEnable bool) // Init initializes the algorithm with the given bin and options.
	packing() *PackResult                     // Pack packs the rectangles into the bin.
}

// basicAlgo basicAlgorithms
type basicAlgo struct {
	bin         *shape.Bin // The bin to pack
	allowRotate bool       // Whether to allow rectangle rotation
}

func (algo *basicAlgo) init(bin *shape.Bin, isRotateEnable bool) {
	algo.bin = bin
	algo.allowRotate = isRotateEnable
}

func (algo *basicAlgo) packing() *PackResult {
	var packedRects []*shape.RectPacked
	var unpackedRects []*shape.Rect
	totalArea := 0
	currentX, currentY := 0, 0
	maxYInRow := 0
	for _, reqRect := range algo.bin.ReqPackRectList {
		if currentX+reqRect.GetW() > algo.bin.GetW() {
			currentX = 0
			currentY += maxYInRow
			maxYInRow = 0
		}
		canPlace := false
		placed := shape.NewRectPacked(*reqRect, 0, 0, false)
		if currentX+reqRect.GetW() <= algo.bin.GetW() && currentY+reqRect.GetH() <= algo.bin.GetH() {
			placed.X = currentX
			placed.Y = currentY
			canPlace = true
		} else if algo.allowRotate && currentX+reqRect.GetH() <= algo.bin.GetW() && currentY+reqRect.GetW() <= algo.bin.GetH() {
			placed.Rotate()
			placed.X = currentX
			placed.Y = currentY
			placed.IsRotated = true
			canPlace = true
		}
		if canPlace {
			packedRects = append(packedRects, placed)
			totalArea += placed.Area()
			currentX += placed.GetW()
			if placed.GetH() > maxYInRow {
				maxYInRow = placed.GetH()
			}
		} else {
			unpackedRects = append(unpackedRects, reqRect)
		}
	}
	fillRate := float64(totalArea) / float64(algo.bin.Area())
	return &PackResult{
		PackedRects:   packedRects,
		UnpackedRects: unpackedRects,
		TotalArea:     totalArea,
		FillRate:      fillRate,
	}
}