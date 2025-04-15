package pack

import (
	"container/heap"
	"fmt"
	"github.com/91xusir/spritepacker/model"
)

// skyline represents a skyline.
type skyline struct{ x, y, len int }

// skylineHeap represents a heap of skylines.
type skylineHeap []skyline

func (h *skylineHeap) Len() int { return len(*h) }

func (h *skylineHeap) Less(i, j int) bool {
	if (*h)[i].y == (*h)[j].y {
		return (*h)[i].x < (*h)[j].x
	}
	return (*h)[i].y < (*h)[j].y
}

func (h *skylineHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *skylineHeap) Push(x any) { *h = append(*h, x.(skyline)) }

func (h *skylineHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// algoSkyline is a packing algo that uses the skyline method.
type algoSkyline struct {
	algoBasic
	skyLineQueue skylineHeap // The smallest pile of skyline
}

func (algo *algoSkyline) init(opt *Options) {
	algo.algoBasic.init(opt)
	algo.skyLineQueue = skylineHeap{}
	heap.Init(&algo.skyLineQueue)
	algo.skyLineQueue = append(algo.skyLineQueue, skyline{x: 0, y: 0, len: algo.w})
}

func (algo *algoSkyline) packing(reqRects []model.Rect) ([]model.Rect, []model.Rect) {
	totalArea := 0
	packedRects := make([]model.Rect, 0, len(reqRects))
	unpackedRects := make([]model.Rect, 0)
	used := make([]bool, len(reqRects))
	for algo.skyLineQueue.Len() != 0 && len(packedRects) < len(reqRects) {
		skyLine := heap.Pop(&algo.skyLineQueue).(skyline)
		hl, hr := algo.getHLHR(skyLine)
		maxRectIndex, maxScore, isRotate := algo.selectMaxScoreRect(skyLine, hl, hr, used, reqRects)
		if maxScore >= 0 {
			packedRect := algo.placeRect(maxRectIndex, skyLine, isRotate, hl, hr, maxScore, reqRects)
			packedRects = append(packedRects, packedRect)
			used[maxRectIndex] = true
			totalArea += reqRects[maxRectIndex].Area()
		} else {
			algo.combineSkylines(skyLine)
		}
	}
	for i, usedFlag := range used {
		if !usedFlag {
			unpackedRects = append(unpackedRects, reqRects[i])
		}
	}
	return packedRects, unpackedRects
}

func (algo *algoSkyline) reset(w, h int) {
	algo.w, algo.h = w, h
	algo.skyLineQueue = skylineHeap{}
	heap.Init(&algo.skyLineQueue)
	algo.skyLineQueue = append(algo.skyLineQueue, skyline{x: 0, y: 0, len: algo.w})
}

// getHLHR selects the highest and lowest skyline
func (algo *algoSkyline) getHLHR(skyLine skyline) (int, int) {
	hl := algo.h - skyLine.y
	hr := algo.h - skyLine.y
	count := 0
	for _, line := range algo.skyLineQueue {
		if line.x+line.len == skyLine.x {
			hl = line.y - skyLine.y
			count++
		} else if line.x == skyLine.x+skyLine.len {
			hr = line.y - skyLine.y
			count++
		}
		if count == 2 {
			break
		}
	}
	return hl, hr
}

// selectMaxScoreRect selects the rectangle with the highest score
func (algo *algoSkyline) selectMaxScoreRect(skyLine skyline, hl, hr int, used []bool, reqRects []model.Rect) (int, int, bool) {
	maxRectIndex, maxScore := -1, -1
	isRotate := false
	for i := range reqRects {
		if !used[i] {
			score := algo.score(reqRects[i].W, reqRects[i].H, skyLine, hl, hr)
			if score > maxScore {
				maxScore = score
				maxRectIndex = i
				isRotate = false
			}
			if algo.allowRotate {
				rotateScore := algo.score(reqRects[i].H, reqRects[i].W, skyLine, hl, hr)
				if rotateScore > maxScore {
					maxScore = rotateScore
					maxRectIndex = i
					isRotate = true
				}
			}
		}
	}
	return maxRectIndex, maxScore, isRotate
}

// placeRect places the rectangle at the specified position
func (algo *algoSkyline) placeRect(maxRectIndex int, skyLine skyline, isRotate bool, hl, hr, maxScore int, reqRects []model.Rect) model.Rect {
	if (hl >= hr && maxScore == 2) || (!(hl >= hr) && (maxScore == 4 || maxScore == 0)) {
		return algo.placeRight(reqRects[maxRectIndex], skyLine, isRotate)
	}
	return algo.placeLeft(reqRects[maxRectIndex], skyLine, isRotate)
}

// placeLeft Place the rectangle to the left
func (algo *algoSkyline) placeLeft(rect model.Rect, skyLine skyline, isRotate bool) model.Rect {
	packedRect := rect.CloneWithPos(skyLine.x, skyLine.y)
	if isRotate {
		packedRect = packedRect.Rotated()
	}
	algo.addSkyLineInQueue(skyLine.x, skyLine.y+packedRect.H, packedRect.W)
	algo.addSkyLineInQueue(skyLine.x+packedRect.W, skyLine.y, skyLine.len-packedRect.W)
	return packedRect
}

// placeRight Place the rectangle to the right
func (algo *algoSkyline) placeRight(rect model.Rect, skyLine skyline, isRotate bool) model.Rect {
	var packedRect model.Rect
	if !isRotate {
		packedRect = rect.CloneWithPos(skyLine.x+skyLine.len-rect.W, skyLine.y)
	} else {
		packedRect = rect.CloneWithPos(skyLine.x+skyLine.len-rect.H, skyLine.y).Rotated()
	}
	algo.addSkyLineInQueue(skyLine.x, skyLine.y, skyLine.len-packedRect.W)
	algo.addSkyLineInQueue(packedRect.X, skyLine.y+packedRect.H, packedRect.W)
	return packedRect
}

// addSkyLineInQueue adds a skyline to the queue
func (algo *algoSkyline) addSkyLineInQueue(x, y, len int) {
	if len > 0 {
		skyLine := skyline{x: x, y: y, len: len}
		heap.Push(&algo.skyLineQueue, skyLine)
	}
}

func (algo *algoSkyline) combineSkylines(skyLine skyline) {
	b := false
	for i, line := range algo.skyLineQueue {
		if skyLine.y <= line.y {
			if skyLine.x == line.x+line.len {
				heap.Remove(&algo.skyLineQueue, i)
				b = true
				skyLine.x = line.x
				skyLine.y = line.y
				skyLine.len = line.len + skyLine.len
				break
			}
			if skyLine.x+skyLine.len == line.x {
				heap.Remove(&algo.skyLineQueue, i)
				b = true
				skyLine.y = line.y
				skyLine.len = line.len + skyLine.len
				break
			}
		}
	}
	if b {
		heap.Push(&algo.skyLineQueue, skyLine)
	}
}

// score calculates the score of the rectangle at the specified position
func (algo *algoSkyline) score(w, h int, skyLine skyline, hl, hr int) int {
	if skyLine.len < w || skyLine.y+h > algo.h {
		return -1
	}
	var high, low int
	if hl >= hr {
		high = hl
		low = hr
	} else {
		high = hr
		low = hl
	}
	if w == skyLine.len && h == high {
		return 7
	} else if w == skyLine.len && h == low {
		return 6
	} else if w == skyLine.len && h > high {
		return 5
	} else if w < skyLine.len && h == high {
		return 4
	} else if w == skyLine.len && h < high && h > low {
		return 3
	} else if w < skyLine.len && h == low {
		return 2
	} else if w == skyLine.len && h < low {
		return 1
	} else if w < skyLine.len && h != high {
		return 0
	}
	panic(fmt.Sprintf("w = %d , h = %d , high = %d , low = %d , skyline = %+v", w, h, high, low, skyLine))
}
