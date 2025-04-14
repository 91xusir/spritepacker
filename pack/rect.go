package pack

import (
	"fmt"
	"image"
	"os"
)

// Size represents a size with width and height.
// The width and height must be positive.
type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

func (s Size) Clone() Size {
	return Size{W: s.W, H: s.H}
}

func NewSize(w, h int) Size {
	if w <= 0 || h <= 0 {
		_, _ = fmt.Fprintf(os.Stderr, "size dimensions must be positive, got: w=%d, h=%d\n", w, h)
		w, h = 1, 1
	}
	return Size{W: w, H: h}
}
func (s Size) Area() int {
	return s.W * s.H
}
func (s Size) Rotated() Size {
	return Size{W: s.H, H: s.W}
}

func (s Size) PowerOfTwo() Size {
	s.W = NextPowerOfTwo(s.W)
	s.H = NextPowerOfTwo(s.H)
	return s
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewPoint(x, y int) Point {
	if x < 0 || y < 0 {
		_, _ = fmt.Fprintf(os.Stderr, "point coordinates must be positive, got: x=%d, y=%d\n", x, y)
		x, y = 0, 0
	}
	return Point{X: x, Y: y}
}

// Rect represents an immutable rectangle using value semantics.
type Rect struct {
	Point
	Size
	Id        int  `json:"-"`
	IsRotated bool `json:"rotated,omitempty"`
}

// NewRect creates a new Rect value.
func NewRect(x, y, w, h, id int) Rect {
	return Rect{
		Point:     NewPoint(x, y),
		Size:      NewSize(w, h),
		Id:        id,
		IsRotated: false,
	}
}
func NewRectBySize(w, h int) Rect {
	return NewRect(0, 0, w, h, 0)
}
func NewRectBySizeAndId(w, h, id int) Rect {
	return NewRect(0, 0, w, h, id)
}
func NewRectByPosAndSize(x, y, w, h int) Rect {
	return NewRect(x, y, w, h, 0)
}
func (r Rect) isContainedIn(rect Rect) bool {
	return r.X >= rect.X && r.Y >= rect.Y && r.X+r.W <= rect.X+rect.W && r.Y+r.H <= rect.Y+rect.H
}

func (r Rect) CloneWithPos(x, y int) Rect {
	return Rect{
		Point:     NewPoint(x, y),
		Size:      r.Size,
		Id:        r.Id,
		IsRotated: r.IsRotated,
	}
}
func (r Rect) CloneWithSize(w, h int) Rect {
	return Rect{
		Point:     r.Point,
		Size:      NewSize(w, h),
		Id:        r.Id,
		IsRotated: r.IsRotated,
	}
}

func (r Rect) Clone() Rect {
	return Rect{
		Point:     r.Point,
		Size:      r.Size,
		Id:        r.Id,
		IsRotated: r.IsRotated,
	}
}

func (r Rect) Rotated() Rect {
	return Rect{
		Point:     r.Point,
		Size:      r.Size.Rotated(),
		Id:        r.Id,
		IsRotated: !r.IsRotated,
	}
}

func (r Rect) ToImageRect() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.W, r.Y+r.H)
}

func (r Rect) String() string {
	return fmt.Sprintf("[Id: %d, X: %d, Y: %d, W: %d, H: %d, IsRotated: %v]", r.Id, r.X, r.Y, r.W, r.H, r.IsRotated)
}

// Bin represents a Bin with a width, height,and a list of rectangles have been packed.
type Bin struct {
	Size
	PackedRects []Rect
	UsedArea    int
}

// NewBin creates a new bin with the given width, height, and list of rectangles to be packed.
func NewBin(w, h int, req []Rect) Bin {
	bin := Bin{
		Size:        NewSize(w, h),
		PackedRects: req,
	}
	return bin
}
func (b Bin) FillRate() float64 {
	return float64(b.UsedArea) / float64(b.Area())
}

func (b Bin) String() string {
	return fmt.Sprintf("Bin{W:%d, H:%d,UsedArea:%d, FillRate:%.2f%%}", b.W, b.H, b.UsedArea, b.FillRate()*100)
}

func addPadding(rect *Rect, padding int) {
	rect.W += padding
	rect.H += padding
}

func removePadding(rect *Rect, padding int) {
	rect.W -= padding
	rect.H -= padding
}
