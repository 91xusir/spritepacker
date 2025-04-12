package pack

import (
	"container/heap"
	"fmt"
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

func (p *algoSkyline) init(opt *Options) {
	p.algoBasic.init(opt)
	p.skyLineQueue = skylineHeap{}
	heap.Init(&p.skyLineQueue)
	p.skyLineQueue = append(p.skyLineQueue, skyline{x: 0, y: 0, len: p.w})
}

func (p *algoSkyline) packing(reqRects []Rect) PackedResult {
	totalArea := 0
	packedRectList := make([]PackedRect, 0, len(reqRects))
	used := make([]bool, len(reqRects))
	for p.skyLineQueue.Len() != 0 && len(packedRectList) < len(reqRects) {
		skyLine := heap.Pop(&p.skyLineQueue).(skyline)
		hl, hr := p.getHLHR(skyLine)
		maxRectIndex, maxScore, isRotate := p.selectMaxScoreRect(skyLine, hl, hr, used, reqRects)
		if maxScore >= 0 {
			packedRect := p.placeRect(maxRectIndex, skyLine, isRotate, hl, hr, maxScore, reqRects)
			packedRectList = append(packedRectList, packedRect)
			used[maxRectIndex] = true
			totalArea += reqRects[maxRectIndex].W * reqRects[maxRectIndex].H
		} else {
			p.combineSkylines(skyLine)
		}
	}
	unpackedRects := make([]Rect, 0)
	for i, usedFlag := range used {
		if !usedFlag {
			unpackedRects = append(unpackedRects, reqRects[i])
		}
	}
	fillRate := float64(totalArea) / float64(p.w*p.h)
	bin := NewBin(p.w, p.h, packedRectList, totalArea, fillRate)

	return PackedResult{
		Bin:           bin,
		UnpackedRects: unpackedRects,
	}
}

// getHLHR selects the highest and lowest skyline
func (p *algoSkyline) getHLHR(skyLine skyline) (int, int) {
	hl := p.h - skyLine.y
	hr := p.h - skyLine.y
	count := 0
	for _, line := range p.skyLineQueue {
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
func (p *algoSkyline) selectMaxScoreRect(skyLine skyline, hl, hr int, used []bool, reqRects []Rect) (int, int, bool) {
	maxRectIndex, maxScore := -1, -1
	isRotate := false
	for i := range reqRects {
		if !used[i] {
			score := p.score(reqRects[i].W, reqRects[i].H, skyLine, hl, hr)
			if score > maxScore {
				maxScore = score
				maxRectIndex = i
				isRotate = false
			}
			if p.allowRotate {
				rotateScore := p.score(reqRects[i].H, reqRects[i].W, skyLine, hl, hr)
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
func (p *algoSkyline) placeRect(maxRectIndex int, skyLine skyline, isRotate bool, hl, hr, maxScore int, reqRects []Rect) PackedRect {
	if (hl >= hr && maxScore == 2) || (!(hl >= hr) && (maxScore == 4 || maxScore == 0)) {
		return p.placeRight(reqRects[maxRectIndex], skyLine, isRotate)
	}
	return p.placeLeft(reqRects[maxRectIndex], skyLine, isRotate)
}

// placeLeft Place the rectangle to the left

func (p *algoSkyline) placeLeft(rect Rect, skyLine skyline, isRotate bool) PackedRect {
	packedRect := NewRectPacked(skyLine.x, skyLine.y, rect)
	if isRotate {
		packedRect = packedRect.Rotated()
	}
	p.addSkyLineInQueue(skyLine.x, skyLine.y+packedRect.H, packedRect.W)
	p.addSkyLineInQueue(skyLine.x+packedRect.W, skyLine.y, skyLine.len-packedRect.W)
	return packedRect
}

// placeRight Place the rectangle to the right

func (p *algoSkyline) placeRight(rect Rect, skyLine skyline, isRotate bool) PackedRect {
	var packedRect PackedRect
	if !isRotate {
		packedRect = NewRectPacked(skyLine.x+skyLine.len-rect.W, skyLine.y, rect)
	} else {
		packedRect = NewRectPacked(skyLine.x+skyLine.len-rect.H, skyLine.y, rect).Rotated()
	}
	p.addSkyLineInQueue(skyLine.x, skyLine.y, skyLine.len-packedRect.W)
	p.addSkyLineInQueue(packedRect.X, skyLine.y+packedRect.H, packedRect.W)
	return packedRect
}

// addSkyLineInQueue adds a skyline to the queue
func (p *algoSkyline) addSkyLineInQueue(x, y, len int) {
	if len > 0 {
		skyLine := skyline{x: x, y: y, len: len}
		heap.Push(&p.skyLineQueue, skyLine)
	}
}

func (p *algoSkyline) combineSkylines(skyLine skyline) {
	b := false
	for i, line := range p.skyLineQueue {
		if skyLine.y <= line.y {
			if skyLine.x == line.x+line.len {
				heap.Remove(&p.skyLineQueue, i)
				b = true
				skyLine.x = line.x
				skyLine.y = line.y
				skyLine.len = line.len + skyLine.len
				break
			}
			if skyLine.x+skyLine.len == line.x {
				heap.Remove(&p.skyLineQueue, i)
				b = true
				skyLine.y = line.y
				skyLine.len = line.len + skyLine.len
				break
			}
		}
	}
	if b {
		heap.Push(&p.skyLineQueue, skyLine)
	}
}

// score calculates the score of the rectangle at the specified position
func (p *algoSkyline) score(w, h int, skyLine skyline, hl, hr int) int {
	if skyLine.len < w || skyLine.y+h > p.h {
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
