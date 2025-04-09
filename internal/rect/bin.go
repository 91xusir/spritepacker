package rect

// Bin represents a bin with a width, height, and a list of rectangles.
// It is used to store the dimensions of the bin and the rectangles that need to be packed into it.
// If IsRotateEnable is true, the rectangles can be rotated to fit into the bin.
// If IsRotateEnable is false, the rectangles cannot be rotated.

type Bin struct {
	// The width and height of the bin
	W, H int
	// Rectangle list
	ReqPackRectList []*Rect
	// Whether to allow rectangle rotation
	IsRotateEnable bool
}
