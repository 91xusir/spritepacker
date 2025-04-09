package spritepacker

import "rectpack2d/spritepacker/internal/rect"

// Bin usedToHoldRectangularContainers

type Bin struct {
	// The width of the boundary

	W int
	// The height of the boundary

	H int
	// Rectangle list

	ReqPackRectList []*rect.Rect
	// Whether to allow rectangle rotation

	IsRotateEnable bool
}
