package pack

import (
	"math"
)

type FreeRect struct {
	Rect
	X, Y int // 左上角坐标
}

// algoMaxrects 装箱器实现
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

// packing 打包
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

// insert 插入单个矩形
func (algo *algoMaxrects) insert(rect Rect) (PackedRect, bool) {
	bestNode := algo.findBestPosition(rect)
	if bestNode.H == 0 {
		return PackedRect{}, false
	}
	algo.placeRect(bestNode)
	return bestNode, true
}

// findBestPosition 查找最佳放置位置
func (algo *algoMaxrects) findBestPosition(rect Rect) PackedRect {
	var bestNode PackedRect
	bestScore := math.MaxInt
	for _, freeRect := range algo.freeRects {
		// 不旋转
		if freeRect.W >= rect.W && freeRect.H >= rect.H {
			score := algo.calculateScore(freeRect, rect.W, rect.H)
			if score < bestScore {
				bestNode = NewRectPacked(rect, freeRect.X, freeRect.Y)
				bestScore = score
			}
		}

		// 尝试旋转
		if algo.allowRotate && freeRect.W >= rect.H && freeRect.H >= rect.W {
			score := algo.calculateScore(freeRect, rect.H, rect.W)
			if score < bestScore {
				bestNode = NewRectPacked(rect, freeRect.X, freeRect.Y).Rotated()
				bestScore = score
			}
		}
	}
	return bestNode
}

// calculateScore 根据策略计算放置分数
func (algo *algoMaxrects) calculateScore(freeRect FreeRect, rectW, rectH int) int {
	switch algo.method {
	case BestShortSideFit:
		return minInt(freeRect.W-rectW, freeRect.H-rectH)
	case BestLongSideFit:
		return maxInt(freeRect.W-rectW, freeRect.H-rectH)
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

// calculateContactPoint 计算接触点数
func (algo *algoMaxrects) calculateContactPoint(freeRect FreeRect, rectW, rectH int) int {
	contactScore := 0
	// 预计算新矩形的边界
	newRect := FreeRect{
		X:    freeRect.X,
		Y:    freeRect.Y,
		Rect: Rect{H: rectH, W: rectW},
	}

	// 检查与容器边界的接触
	if newRect.X == 0 || newRect.X+newRect.W == algo.w {
		contactScore += newRect.H
	}
	if newRect.Y == 0 || newRect.Y+newRect.H == algo.h {
		contactScore += newRect.W
	}

	// 检查与已放置矩形的接触
	for _, usedRect := range algo.usedRects {
		if newRect.X == usedRect.X+usedRect.W || newRect.X+newRect.W == usedRect.X {
			// 垂直接触
			overlap := minInt(newRect.Y+newRect.H, usedRect.Y+usedRect.H) -
				maxInt(newRect.Y, usedRect.Y)
			if overlap > 0 {
				contactScore += overlap
			}
		}

		if newRect.Y == usedRect.Y+usedRect.H || newRect.Y+newRect.H == usedRect.Y {
			// 水平接触
			overlap := minInt(newRect.X+newRect.W, usedRect.X+usedRect.W) -
				maxInt(newRect.X, usedRect.X)
			if overlap > 0 {
				contactScore += overlap
			}
		}
	}
	return contactScore
}

// placeRect 放置矩形并分割剩余空间
func (algo *algoMaxrects) placeRect(rect PackedRect) {
	// 分割所有相交的空闲矩形
	for i := 0; i < len(algo.freeRects); {
		if algo.splitFreeRect(algo.freeRects[i], rect) {
			algo.freeRects = append(algo.freeRects[:i], algo.freeRects[i+1:]...)
		} else {
			i++
		}
	}
	// 清理无效的空闲矩形
	algo.pruneFreeList()
	algo.usedRects = append(algo.usedRects, rect)
}

// splitFreeRect 分割空闲矩形
func (algo *algoMaxrects) splitFreeRect(freeRect FreeRect, usedRect PackedRect) bool {
	// 检查是否有重叠
	if usedRect.X >= freeRect.X+freeRect.W || usedRect.X+usedRect.W <= freeRect.X ||
		usedRect.Y >= freeRect.Y+freeRect.H || usedRect.Y+usedRect.H <= freeRect.Y {
		return false
	}

	// 生成新的空闲矩形（上下左右四个部分）
	// 上部分
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

	// 下部分
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

	// 左部分
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

	// 右部分
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

// pruneFreeList 清理无效的空闲矩形
func (algo *algoMaxrects) pruneFreeList() {
	// 移除被完全包含的空闲矩形
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

// isContained 检查rect1是否完全包含在rect2中
func (algo *algoMaxrects) isContained(rect1, rect2 FreeRect) bool {
	return rect1.X >= rect2.X && rect1.Y >= rect2.Y &&
		rect1.X+rect1.W <= rect2.X+rect2.W &&
		rect1.Y+rect1.H <= rect2.Y+rect2.H
}
