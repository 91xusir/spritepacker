package pack

import (
	"math"
)

type FreeRect struct {
	Rect
	X, Y int
}

type algoMaxrects struct {
	algoBasic
	usedRects []PackedRect
	freeRects []FreeRect
	method    Heuristic
}

func (algo *algoMaxrects) init(opt *Options) {
	algo.algoBasic.init(opt)
	algo.freeRects = []FreeRect{{Rect: Rect{W: algo.w, H: algo.h}, X: 0, Y: 0}}
	algo.method = opt.heuristic
}

func (algo *algoMaxrects) reset(w, h int) {
	algo.algoBasic.reset(w, h)
	algo.freeRects = []FreeRect{{Rect: Rect{W: algo.w, H: algo.h}, X: 0, Y: 0}}
}

func (algo *algoMaxrects) packing(reqRects []Rect) PackedResult {
	result := PackedResult{
		UnpackedRects: make([]Rect, 0),
		Bin:           NewBin(algo.w, algo.h, make([]PackedRect, 0), 0, 0),
	}
	for _, rect := range reqRects {

		if packedRect, ok := algo.insert(rect); ok {
			result.Bin.UsedArea += rect.Area()
			result.Bin.PackedRects = append(result.Bin.PackedRects, packedRect)
		} else {
			result.UnpackedRects = append(result.UnpackedRects, rect)
		}
	}
	if len(result.Bin.PackedRects) > 0 {
		result.Bin.FillRate = float64(result.Bin.UsedArea) / float64(algo.w*algo.h)
	}
	return result
}

func (algo *algoMaxrects) insert(rect Rect) (PackedRect, bool) {
	bestNode := algo.findBestPosition(rect)
	if bestNode.H == 0 {
		return PackedRect{}, false
	}
	algo.placeRect(bestNode)
	return bestNode, true
}

func (algo *algoMaxrects) findBestPosition(rect Rect) PackedRect {
	var bestNode PackedRect
	bestScore := math.MaxInt
	for _, freeRect := range algo.freeRects {
		if freeRect.W >= rect.W && freeRect.H >= rect.H {
			score := algo.calculateScore(freeRect, rect.W, rect.H)
			if score < bestScore {
				bestNode = NewRectPacked(freeRect.X, freeRect.Y, rect)
				bestScore = score
			}
		}
		if algo.allowRotate && freeRect.W >= rect.H && freeRect.H >= rect.W {
			score := algo.calculateScore(freeRect, rect.H, rect.W)
			if score < bestScore {
				bestNode = NewRectPacked(freeRect.X, freeRect.Y, rect).Rotated()
				bestScore = score
			}
		}
	}
	return bestNode
}

func (algo *algoMaxrects) calculateScore(freeRect FreeRect, rectW, rectH int) int {
	switch algo.method {
	case BestShortSideFit:
		return MinInt(freeRect.W-rectW, freeRect.H-rectH)
	case BestLongSideFit:
		return MaxInt(freeRect.W-rectW, freeRect.H-rectH)
	case BestAreaFit:
		return freeRect.W*freeRect.H - rectW*rectH
	case BottomLeftFit:
		return freeRect.Y + rectH
	case ContactPointFit:
		return -algo.calculateContactPoint(freeRect, rectW, rectH)
	default:
		return 0
	}
}

func (algo *algoMaxrects) calculateContactPoint(freeRect FreeRect, rectW, rectH int) int {
	contactScore := 0
	newRect := FreeRect{
		X:    freeRect.X,
		Y:    freeRect.Y,
		Rect: Rect{H: rectH, W: rectW},
	}
	if newRect.X == 0 || newRect.X+newRect.W == algo.w {
		contactScore += newRect.H
	}
	if newRect.Y == 0 || newRect.Y+newRect.H == algo.h {
		contactScore += newRect.W
	}
	for _, usedRect := range algo.usedRects {
		if newRect.X == usedRect.X+usedRect.W || newRect.X+newRect.W == usedRect.X {
			overlap := MinInt(newRect.Y+newRect.H, usedRect.Y+usedRect.H) -
				MaxInt(newRect.Y, usedRect.Y)
			if overlap > 0 {
				contactScore += overlap
			}
		}

		if newRect.Y == usedRect.Y+usedRect.H || newRect.Y+newRect.H == usedRect.Y {
			overlap := MinInt(newRect.X+newRect.W, usedRect.X+usedRect.W) -
				MaxInt(newRect.X, usedRect.X)
			if overlap > 0 {
				contactScore += overlap
			}
		}
	}
	return contactScore
}

func (algo *algoMaxrects) placeRect(rect PackedRect) {
	for i := 0; i < len(algo.freeRects); {
		if algo.splitFreeRect(algo.freeRects[i], rect) {
			algo.freeRects = append(algo.freeRects[:i], algo.freeRects[i+1:]...)
		} else {
			i++
		}
	}
	algo.pruneFreeList()
	algo.usedRects = append(algo.usedRects, rect)
}

func (algo *algoMaxrects) splitFreeRect(freeRect FreeRect, usedRect PackedRect) bool {
	if usedRect.X >= freeRect.X+freeRect.W || usedRect.X+usedRect.W <= freeRect.X ||
		usedRect.Y >= freeRect.Y+freeRect.H || usedRect.Y+usedRect.H <= freeRect.Y {
		return false
	}

	// upper part
	if usedRect.Y > freeRect.Y {
		algo.freeRects = append(algo.freeRects, FreeRect{
			X: freeRect.X,
			Y: freeRect.Y,
			Rect: Rect{
				W: freeRect.W,
				H: usedRect.Y - freeRect.Y,
			},
		})
	}

	// lower part
	if usedRect.Y+usedRect.H < freeRect.Y+freeRect.H {
		algo.freeRects = append(algo.freeRects, FreeRect{
			X: freeRect.X,
			Y: usedRect.Y + usedRect.H,
			Rect: Rect{
				W: freeRect.W,
				H: freeRect.Y + freeRect.H - (usedRect.Y + usedRect.H),
			},
		})
	}

	// left part
	if usedRect.X > freeRect.X {
		algo.freeRects = append(algo.freeRects, FreeRect{
			X: freeRect.X,
			Y: freeRect.Y,
			Rect: Rect{
				W: usedRect.X - freeRect.X,
				H: freeRect.H,
			},
		})
	}

	// right part
	if usedRect.X+usedRect.W < freeRect.X+freeRect.W {
		algo.freeRects = append(algo.freeRects, FreeRect{
			X: usedRect.X + usedRect.W,
			Y: freeRect.Y,
			Rect: Rect{
				W: freeRect.X + freeRect.W - (usedRect.X + usedRect.W),
				H: freeRect.H,
			},
		})
	}
	return true
}

func (algo *algoMaxrects) pruneFreeList() {
	for i := 0; i < len(algo.freeRects); i++ {
		for j := i + 1; j < len(algo.freeRects); {
			if algo.isContained(algo.freeRects[i], algo.freeRects[j]) {
				algo.freeRects = append(algo.freeRects[:i], algo.freeRects[i+1:]...)
				i--
				break
			}
			if algo.isContained(algo.freeRects[j], algo.freeRects[i]) {
				algo.freeRects = append(algo.freeRects[:j], algo.freeRects[j+1:]...)
			} else {
				j++
			}
		}
	}
}

func (algo *algoMaxrects) isContained(rect1, rect2 FreeRect) bool {
	return rect1.X >= rect2.X && rect1.Y >= rect2.Y &&
		rect1.X+rect1.W <= rect2.X+rect2.W &&
		rect1.Y+rect1.H <= rect2.Y+rect2.H
}
