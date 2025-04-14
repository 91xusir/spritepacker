package pack

// Algorithm defines the packing algorithm used.
type Algorithm int

const (
	AlgoBasic Algorithm = iota
	AlgoSkyline
	AlgoMaxRects
	MaxAlgoIndex
)

// Heuristic defines the heuristic used by the MaxRects algorithm.
type Heuristic int

const (
	BestShortSideFit Heuristic = iota
	BestLongSideFit
	BestAreaFit
	BottomLeftFit
	ContactPointFit
	MaxHeuristicsIndex
)

// algo is the interface that wraps the Pack method.
type algo interface {
	init(opt *Options)                        // Init initializes the algo with the given bin and Options.
	packing(reqRects []Rect) ([]Rect, []Rect) // Pack packs the given rectangles into the bin and returns Packed and Unpacked rectangles.
	reset(w, h int)                           // ResetWH resets the width and height of the bin.
}

// algoBasic basicAlgorithms
type algoBasic struct {
	w, h        int  // The width and height of the bin
	allowRotate bool // Whether to allow rectangle rotation
}

func (algo *algoBasic) reset(w, h int) {
	algo.w, algo.h = w, h
}

func (algo *algoBasic) init(opt *Options) {
	algo.w = opt.maxW
	algo.h = opt.maxH
	algo.allowRotate = opt.allowRotate
}

func (algo *algoBasic) packing(reqRects []Rect) ([]Rect, []Rect) {
	packedRects := make([]Rect, 0, len(reqRects))
	unpackedRects := make([]Rect, 0)
	totalArea := 0
	currentX, currentY := 0, 0
	maxYInRow := 0
	for _, reqRect := range reqRects {
		if currentX+reqRect.W > algo.w {
			currentX = 0
			currentY += maxYInRow
			maxYInRow = 0
		}
		canPlace := false
		placed := reqRect.Clone()
		if currentX+reqRect.W <= algo.w && currentY+reqRect.H <= algo.h {
			placed.X = currentX
			placed.Y = currentY
			canPlace = true
		} else if algo.allowRotate && currentX+reqRect.H <= algo.w && currentY+reqRect.W <= algo.h {
			//placed.IsRotated = true
			placed = placed.Rotated()
			placed.X = currentX
			placed.Y = currentY
			canPlace = true
		}
		if canPlace {
			currentX += placed.W
			if placed.H > maxYInRow {
				maxYInRow = placed.H
			}
			totalArea += placed.Area()
			packedRects = append(packedRects, placed)
		} else {
			unpackedRects = append(unpackedRects, reqRect)
		}
	}
	return packedRects, unpackedRects
}
