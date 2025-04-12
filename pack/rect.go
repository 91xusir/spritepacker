package pack

import "fmt"

// Rect represents an immutable rectangle using value semantics.
type Rect struct {
	W  int // Width (must be > 0)
	H  int // Height (must be > 0)
	Id int // Optional identifier
}

// NewRect creates a new Rect value.
func NewRect(w, h int) Rect {
	return NewRectById(w, h, 0)
}

// NewRectById creates a configured Rect value.
func NewRectById(w, h int, id int) Rect {
	if w <= 0 || h <= 0 {
		fmt.Printf("rect dimensions must be positive, got: w=%d, h=%d\n", w, h)
		w, h = 1, 1
	}
	return Rect{W: w, H: h, Id: id}
}

// ----- Core Methods -----

// Area returns the area of the rectangle.
//
//go:inline
func (r Rect) Area() int { return r.W * r.H }

// Ratio returns the width to height ratio.
func (r Rect) Ratio() float64 { return float64(r.W) / float64(r.H) }

// MaxSide returns the longer side length.
func (r Rect) MaxSide() int { return maxInt(r.W, r.H) }

// MinSide returns the shorter side length.
func (r Rect) MinSide() int { return minInt(r.W, r.H) }

// Perimeter returns the perimeter length.
func (r Rect) Perimeter() int { return (r.W + r.H) * 2 }

// ----- Builder Pattern -----

// WithW returns a new Rect with updated width.
func (r Rect) WithSize(w, h int) Rect {
	if w <= 0 || h <= 0 {
		fmt.Printf("rect dimensions must be positive, got: w=%d, h=%d\n", w, h)
		w, h = 1, 1
	}
	return Rect{W: w, H: h, Id: r.Id}
}

// WithId returns a new Rect with updated ID.
func (r Rect) WithId(id int) Rect {
	return Rect{W: r.W, H: r.H, Id: id}
}

// ----- Rotation -----

// Rotated returns a new rotated rectangle.
func (r Rect) Rotated() Rect {
	return Rect{W: r.H, H: r.W, Id: r.Id}
}

// ----- Display -----

func (r Rect) String() string {
	return fmt.Sprintf("Rect{ID: %-5v W: %-4d H: %-4d}",
		r.Id, r.W, r.H)
}

// PackedRect represents a packed rectangle
type PackedRect struct {
	Rect
	X, Y      int // left-top point
	IsRotated bool
}

// NewRectPacked creates a new PackedRect with the given ID, x, y, width, height, and rotation flag.
// If isRotated is true, the rectangle is rotated 90 degrees.
func NewRectPacked(rect Rect, x, y int) PackedRect {
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
