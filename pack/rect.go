package pack

import (
	"fmt"
	"image"
)

// Rect represents an immutable rectangle using value semantics.
type Rect struct {
	W  int
	H  int
	Id int
}

// NewRect creates a new Rect value.
func NewRect(w, h int) Rect {
	return NewRectById(w, h, 0)
}

// NewRectById creates a configured Rect value.
func NewRectById(w, h int, id int) Rect {
	if w <= 0 || h <= 0 {
		fmt.Printf("rect dimensions must be positive, got: id=%d w=%d, h=%d\n", id, w, h)
		w, h = 1, 1
	}
	return Rect{W: w, H: h, Id: id}
}

// Area returns the area of the rectangle.
//
//go:inline
func (r Rect) Area() int { return r.W * r.H }

// Rotated returns a new rotated rectangle.
func (r Rect) Rotated() Rect {
	return Rect{W: r.H, H: r.W, Id: r.Id}
}

func (r Rect) String() string {
	return fmt.Sprintf("Rect{ID: %-5v W: %-4d H: %-4d}",
		r.Id, r.W, r.H)
}

// PackedRect represents a packed rectangle
type PackedRect struct {
	X int
	Y int
	Rect
	IsRotated bool
}

// NewRectPacked creates a new PackedRect with the given ID, x, y, width, height, and rotation flag.
// If isRotated is true, the rectangle is rotated 90 degrees.
func NewRectPacked(x, y int, rect Rect) PackedRect {
	return PackedRect{
		Rect:      rect,
		X:         x,
		Y:         y,
		IsRotated: false,
	}
}

func (r PackedRect) Rotated() PackedRect {
	return PackedRect{
		Rect:      r.Rect.Rotated(),
		X:         r.X,
		Y:         r.Y,
		IsRotated: !r.IsRotated,
	}
}

func (r PackedRect) ToImageRect() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.W, r.Y+r.H)
}
func (r PackedRect) ToRectangle() Rectangle {
	return Rectangle{
		X: r.X,
		Y: r.Y,
		W: r.W,
		H: r.H,
	}
}

func (r PackedRect) String() string {
	return fmt.Sprintf("[Id: %v, X: %d, Y: %d, W: %d, H: %d, IsRotated: %v]", r.Id, r.X, r.Y, r.W, r.H, r.IsRotated)
}

// Bin represents a Bin with a width, height,and a list of rectangles have been packed.
type Bin struct {
	Rect
	PackedRects []PackedRect
	UsedArea    int
	FillRate    float64
}

// NewBin creates a new bin with the given width, height, and list of rectangles to be packed.
func NewBin(w, h int, req []PackedRect, usedArea int, fillRate float64) Bin {
	bin := Bin{
		Rect:        NewRect(w, h),
		PackedRects: req,
		FillRate:    fillRate,
		UsedArea:    usedArea,
	}
	return bin
}

func (b Bin) String() string {
	return fmt.Sprintf("{Id:%v, W:%d, H:%d, FillRate:%.2f%%, PackedRects:%v}", b.Id, b.W, b.H, b.FillRate*100, b.PackedRects)
}

// PackedResult represents the result of packing rectangles into a bin.
type PackedResult struct {
	Bin           Bin
	UnpackedRects []Rect
}

func (r PackedResult) String() string {
	return fmt.Sprintf("PackedResult{Bin:%s, UnpackedRects:%v}", r.Bin, r.UnpackedRects)
}

func addPadding(rect *Rect, padding int) {
	rect.W += padding
	rect.H += padding
}

func removePadding(rect *Rect, padding int) {
	rect.W -= padding
	rect.H -= padding
}
