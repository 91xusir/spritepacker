package pack

import (
	"github.com/91xusir/spritepacker/model"
	"github.com/91xusir/spritepacker/utils"
	"math"
)

type algoMaxrects struct {
	algoBasic
	usedRects []model.Rect
	freeRects []model.Rect
	method    Heuristic
}

func (algo *algoMaxrects) init(opt *Options) {
	algo.algoBasic.init(opt)
	algo.freeRects = []model.Rect{model.NewRectByPosAndSize(0, 0, algo.w, algo.h)}
	algo.usedRects = make([]model.Rect, 0)
	algo.method = opt.heuristic
}

func (algo *algoMaxrects) reset(w, h int) {
	algo.algoBasic.reset(w, h)
	algo.freeRects = []model.Rect{model.NewRectByPosAndSize(0, 0, algo.w, algo.h)}
	algo.usedRects = make([]model.Rect, 0)
}

func (algo *algoMaxrects) packing(reqRects []model.Rect) ([]model.Rect, []model.Rect) {
	packedRects := make([]model.Rect, 0, len(reqRects))
	unpackedRects := make([]model.Rect, 0)
	for _, rect := range reqRects {
		if packedRect, ok := algo.insert(rect); ok {
			packedRects = append(packedRects, packedRect)
		} else {
			unpackedRects = append(unpackedRects, rect)
		}
	}
	return packedRects, unpackedRects
}

func (algo *algoMaxrects) insert(rect model.Rect) (model.Rect, bool) {
	bestNode := algo.findBestPosition(rect)
	if bestNode.H == 0 {
		return model.Rect{}, false
	}
	algo.placeRect(bestNode)
	return bestNode, true
}

func (algo *algoMaxrects) findBestPosition(rect model.Rect) model.Rect {
	var bestNode model.Rect
	bestScore := math.MaxInt
	for _, freeRect := range algo.freeRects {
		if freeRect.W >= rect.W && freeRect.H >= rect.H {
			score := algo.calculateScore(freeRect, rect.W, rect.H)
			if score < bestScore {
				bestNode = rect.CloneWithPos(freeRect.X, freeRect.Y)
				bestScore = score
			}
		}
		if algo.allowRotate && freeRect.W >= rect.H && freeRect.H >= rect.W {
			score := algo.calculateScore(freeRect, rect.H, rect.W)
			if score < bestScore {
				bestNode = rect.CloneWithPos(freeRect.X, freeRect.Y).Rotated()
				bestScore = score
			}
		}
	}
	return bestNode
}

func (algo *algoMaxrects) calculateScore(freeRect model.Rect, rectW, rectH int) int {
	switch algo.method {
	case BestShortSideFit:
		return utils.MinInt(freeRect.W-rectW, freeRect.H-rectH)
	case BestLongSideFit:
		return utils.MaxInt(freeRect.W-rectW, freeRect.H-rectH)
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

func (algo *algoMaxrects) calculateContactPoint(freeRect model.Rect, rectW, rectH int) int {
	contactScore := 0
	//newRect := FreeRect{
	//	X:    freeRect.X,
	//	Y:    freeRect.Y,
	//	Rect: Rect{H: rectH, W: rectW},
	//}
	newRect := freeRect.CloneWithSize(rectW, rectH)

	if newRect.X == 0 || newRect.X+newRect.W == algo.w {
		contactScore += newRect.H
	}
	if newRect.Y == 0 || newRect.Y+newRect.H == algo.h {
		contactScore += newRect.W
	}
	for _, usedRect := range algo.usedRects {
		if newRect.X == usedRect.X+usedRect.W || newRect.X+newRect.W == usedRect.X {
			overlap := utils.MinInt(newRect.Y+newRect.H, usedRect.Y+usedRect.H) -
				utils.MaxInt(newRect.Y, usedRect.Y)
			if overlap > 0 {
				contactScore += overlap
			}
		}

		if newRect.Y == usedRect.Y+usedRect.H || newRect.Y+newRect.H == usedRect.Y {
			overlap := utils.MinInt(newRect.X+newRect.W, usedRect.X+usedRect.W) -
				utils.MaxInt(newRect.X, usedRect.X)
			if overlap > 0 {
				contactScore += overlap
			}
		}
	}
	return contactScore
}

func (algo *algoMaxrects) placeRect(rect model.Rect) {
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

func (algo *algoMaxrects) splitFreeRect(freeRect model.Rect, usedRect model.Rect) bool {
	if usedRect.X >= freeRect.X+freeRect.W || usedRect.X+usedRect.W <= freeRect.X ||
		usedRect.Y >= freeRect.Y+freeRect.H || usedRect.Y+usedRect.H <= freeRect.Y {
		return false
	}

	// upper part
	if usedRect.Y > freeRect.Y {
		/*algo.freeRects = append(algo.freeRects, FreeRect{
			X: freeRect.X,
			Y: freeRect.Y,
			Rect: Rect{
				W: freeRect.W,
				H: usedRect.Y - freeRect.Y,
			},
		})*/
		algo.freeRects = append(algo.freeRects, freeRect.CloneWithSize(freeRect.W, usedRect.Y-freeRect.Y))

	}

	// lower part
	if usedRect.Y+usedRect.H < freeRect.Y+freeRect.H {
		//algo.freeRects = append(algo.freeRects, FreeRect{
		//	X: freeRect.X,
		//	Y: usedRect.Y + usedRect.H,
		//	Rect: Rect{
		//		W: freeRect.W,
		//		H: freeRect.Y + freeRect.H - (usedRect.Y + usedRect.H),
		//	},
		//})
		algo.freeRects = append(algo.freeRects, model.NewRectByPosAndSize(freeRect.X, usedRect.Y+usedRect.H, freeRect.W, freeRect.Y+freeRect.H-(usedRect.Y+usedRect.H)))
	}

	// left part
	if usedRect.X > freeRect.X {

		//algo.freeRects = append(algo.freeRects, FreeRect{
		//	X: freeRect.X,
		//	Y: freeRect.Y,
		//	Rect: Rect{
		//		W: usedRect.X - freeRect.X,
		//		H: freeRect.H,
		//	},
		//})
		algo.freeRects = append(algo.freeRects, freeRect.CloneWithSize(usedRect.X-freeRect.X, freeRect.H))
	}

	// right part
	if usedRect.X+usedRect.W < freeRect.X+freeRect.W {
		//algo.freeRects = append(algo.freeRects, FreeRect{
		//	X: usedRect.X + usedRect.W,
		//	Y: freeRect.Y,
		//	Rect: Rect{
		//		W: freeRect.X + freeRect.W - (usedRect.X + usedRect.W),
		//		H: freeRect.H,
		//	},
		//})
		algo.freeRects = append(algo.freeRects, model.NewRectByPosAndSize(usedRect.X+usedRect.W, freeRect.Y, freeRect.X+freeRect.W-(usedRect.X+usedRect.W), freeRect.H))
	}
	return true
}

func (algo *algoMaxrects) pruneFreeList() {
	for i := 0; i < len(algo.freeRects); i++ {
		for j := i + 1; j < len(algo.freeRects); {
			if algo.freeRects[i].IsContainedIn(algo.freeRects[j]) {
				algo.freeRects = append(algo.freeRects[:i], algo.freeRects[i+1:]...)
				i--
				break
			}
			if algo.freeRects[j].IsContainedIn(algo.freeRects[i]) {
				algo.freeRects = append(algo.freeRects[:j], algo.freeRects[j+1:]...)
			} else {
				j++
			}
		}
	}
}

//func (algo *algoMaxrects) isContained(rect1, rect2 Rect) bool {
//	return rect1.X >= rect2.X && rect1.Y >= rect2.Y &&
//		rect1.X+rect1.W <= rect2.X+rect2.W &&
//		rect1.Y+rect1.H <= rect2.Y+rect2.H
//}
