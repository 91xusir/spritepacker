package packer

import (
	"container/heap"
	"fmt"
	"spritepacker/internal/shape"
)

// skyline represents a skyline.
type skyline struct{ x, y, len int }

// skylineHeap represents a heap of skylines.
type skylineHeap []*skyline

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

func (h *skylineHeap) Push(x any) { *h = append(*h, x.(*skyline)) }

func (h *skylineHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// algoSkyline is a packing algorithm that uses the skyline method.
type algoSkyline struct {
	basicAlgo
	skyLineQueue skylineHeap // The smallest pile of skyline
}

func (algo *algoSkyline) init(bin *shape.Bin, isRotateEnable bool) {
	algo.basicAlgo.init(bin, isRotateEnable)
	algo.skyLineQueue = skylineHeap{}
	heap.Init(&algo.skyLineQueue)
	algo.skyLineQueue = append(algo.skyLineQueue, &skyline{x: 0, y: 0, len: bin.GetW()})
}

func (algo *algoSkyline) packing() *PackResult {
	totalS := 0
	packedRectList := make([]*shape.RectPacked, 0, len(algo.bin.ReqPackRectList))
	used := make([]bool, len(algo.bin.ReqPackRectList))
	for algo.skyLineQueue.Len() != 0 && len(packedRectList) < len(algo.bin.ReqPackRectList) {
		skyLine := heap.Pop(&algo.skyLineQueue).(*skyline)
		hl, hr := algo.getHLHR(skyLine)
		maxRectIndex, maxScore, isRotate := algo.selectMaxScoreRect(skyLine, hl, hr, used)
		if maxScore >= 0 {
			packedRect := algo.placeRect(maxRectIndex, skyLine, isRotate, hl, hr, maxScore)
			packedRectList = append(packedRectList, packedRect)
			used[maxRectIndex] = true
			totalS += algo.bin.ReqPackRectList[maxRectIndex].GetW() * algo.bin.ReqPackRectList[maxRectIndex].GetH()
		} else {
			algo.combineSkylines(skyLine)
		}
	}
	unpackedRects := make([]*shape.Rect, 0)
	for i, usedFlag := range used {
		if !usedFlag {
			unpackedRects = append(unpackedRects, algo.bin.ReqPackRectList[i])
		}
	}
	return &PackResult{
		PackedRects:   packedRectList,
		UnpackedRects: unpackedRects,
		TotalArea:     totalS,
		FillRate:      float64(totalS) / float64(algo.bin.Area()),
	}
}

// getHLHR selects the highest and lowest skyline
func (algo *algoSkyline) getHLHR(skyLine *skyline) (int, int) {
	hl := algo.bin.GetH() - skyLine.y
	hr := algo.bin.GetH() - skyLine.y
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

// selectMaxScsporeRect selects the rectangle with the highest score
func (algo *algoSkyline) selectMaxScoreRect(skyLine *skyline, hl, hr int, used []bool) (int, int, bool) {
	maxRectIndex, maxScore := -1, -1
	isRotate := false
	for i := range algo.bin.ReqPackRectList {
		if !used[i] {
			score := algo.score(algo.bin.ReqPackRectList[i].GetW(), algo.bin.ReqPackRectList[i].GetH(), skyLine, hl, hr)
			if score > maxScore {
				maxScore = score
				maxRectIndex = i
				isRotate = false
			}
			if algo.allowRotate {
				rotateScore := algo.score(algo.bin.ReqPackRectList[i].GetH(), algo.bin.ReqPackRectList[i].GetW(), skyLine, hl, hr)
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
func (algo *algoSkyline) placeRect(maxRectIndex int, skyLine *skyline, isRotate bool, hl, hr, maxScore int) *shape.RectPacked {
	if (hl >= hr && maxScore == 2) || (!(hl >= hr) && (maxScore == 4 || maxScore == 0)) {
		return algo.placeRight(algo.bin.ReqPackRectList[maxRectIndex], skyLine, isRotate)
	}
	return algo.placeLeft(algo.bin.ReqPackRectList[maxRectIndex], skyLine, isRotate)
}

// placeLeft 将矩形靠左放

func (algo *algoSkyline) placeLeft(rect *shape.Rect, skyLine *skyline, isRotate bool) *shape.RectPacked {
	var packedRect *shape.RectPacked
	if !isRotate {
		packedRect = shape.NewRectPacked(*rect, skyLine.x, skyLine.y, isRotate)
	} else {

		packedRect = shape.NewRectPacked(*rect, skyLine.x, skyLine.y, isRotate)
	}
	algo.addSkyLineInQueue(skyLine.x, skyLine.y+packedRect.GetH(), packedRect.GetW())
	algo.addSkyLineInQueue(skyLine.x+packedRect.GetW(), skyLine.y, skyLine.len-packedRect.GetW())
	return packedRect
}

// placeRight Place the rectangle to the right

func (algo *algoSkyline) placeRight(rect *shape.Rect, skyLine *skyline, isRotate bool) *shape.RectPacked {
	var packedRect *shape.RectPacked
	if !isRotate {
		packedRect = shape.NewRectPacked(*rect, skyLine.x+skyLine.len-rect.GetW(), skyLine.y, isRotate)
	} else {
		packedRect = shape.NewRectPacked(*rect, skyLine.x+skyLine.len-rect.GetH(), skyLine.y, isRotate)
	}
	algo.addSkyLineInQueue(skyLine.x, skyLine.y, skyLine.len-packedRect.GetW())
	algo.addSkyLineInQueue(packedRect.X, skyLine.y+packedRect.GetH(), packedRect.GetW())
	return packedRect
}

// addSkyLineInQueue 将指定属性的天际线加入天际线队列

func (algo *algoSkyline) addSkyLineInQueue(x, y, len int) {
	if len > 0 {
		skyLine := &skyline{x: x, y: y, len: len}
		heap.Push(&algo.skyLineQueue, skyLine)
	}
}

func (algo *algoSkyline) combineSkylines(skyLine *skyline) {
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
func (algo *algoSkyline) score(w, h int, skyLine *skyline, hl, hr int) int {
	if skyLine.len < w || skyLine.y+h > algo.bin.GetH() {
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
