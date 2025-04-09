package packer

import "fmt"

// SkyLine represents a skyline.
// x, y: the coordinates of the skyline.
// len: the length of the skyline.
type SkyLine struct {
	x, y, len int
}

// String returns a string representation of the skyline.
// The format is "SkyLine{x=%d, y=%d, len=%d}".
func (s *SkyLine) String() string {
	return fmt.Sprintf("SkyLine{x=%d, y=%d, len=%d}", s.x, s.y, s.len)
}

// NewSkyLine creates a new skyline with the given coordinates and length.
// It returns a pointer to the newly created skyline.
func NewSkyLine(x, y, len int) *SkyLine {
	return &SkyLine{
		x:   x,
		y:   y,
		len: len,
	}
}

// SkyLineHeap represents a heap of skylines.
// It implements the heap.Interface.
type SkyLineHeap []*SkyLine

func (h *SkyLineHeap) Len() int { return len(*h) }

func (h *SkyLineHeap) Less(i, j int) bool {
	if (*h)[i].y == (*h)[j].y {
		return (*h)[i].x < (*h)[j].x
	}
	return (*h)[i].y < (*h)[j].y
}

func (h *SkyLineHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *SkyLineHeap) Push(x any) {
	*h = append(*h, x.(*SkyLine))
}

func (h *SkyLineHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
