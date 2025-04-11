package shape

import "fmt"

// RectInterface is interface for rectangles.
// It provides methods to calculate the area, perimeter, maximum side, minimum side, ratio,
// set width, set height, get width, get height, set ID, and get ID.
// The ID can be any type.
// The width and height must be greater than 0
type RectInterface interface {
	Area() int          // Area returns the area of the rectangle.
	Perimeter() int     // Perimeter returns the perimeter of the rectangle.
	MaxSide() int       // MaxSide returns the maximum side of the rectangle.
	MinSide() int       // MinSide returns the minimum side of the rectangle.
	Ratio() float64     // Ratio returns the ratio of the rectangle.
	SetW(w int) *Rect   // SetW sets the width of the rectangle.
	SetH(h int) *Rect   // SetH sets the height of the rectangle.
	GetW() int          // GetW returns the width of the rectangle.
	GetH() int          // GetH returns the height of the rectangle.
	SetId(id any) *Rect // SetId sets the ID of the rectangle.
	GetId() any         // GetId returns the ID of the rectangle.
	Rotate() *Rect      // Rotate rotates the rectangle.
}

// Rect represents a rectangle.
type Rect struct {
	w, h int // width / height
	id   any // ID of the rectangle, can be any type
}

// NewRect creates a new Rect with the given width and height.
// It returns a pointer to the newly created Rect.
func NewRect(w, h int) *Rect { return NewRectById(w, h, nil) }

// NewRectById creates a new Rect with the given ID, width, and height.
// It returns a pointer to the newly created Rect.
func NewRectById(w int, h int, id any) *Rect {
	return (&Rect{}).SetH(h).SetW(w).SetId(id)
}

func (r *Rect) Area() int      { return r.w * r.h }
func (r *Rect) Ratio() float64 { return float64(r.w) / float64(r.h) }
func (r *Rect) MaxSide() int   { return max(r.w, r.h) }
func (r *Rect) MinSide() int   { return min(r.w, r.h) }
func (r *Rect) Perimeter() int { return (r.w + r.h) << 1 }

func (r *Rect) GetW() int { return r.w }
func (r *Rect) GetH() int { return r.h }

func (r *Rect) SetW(w int) *Rect {
	if w < 0 {
		panic("width must be greater than 0")
	}
	r.w = w
	return r
}

func (r *Rect) SetH(h int) *Rect {
	if h < 0 {
		panic("height must be greater than 0")
	}
	r.h = h
	return r
}

func (r *Rect) SetId(id any) *Rect {
	r.id = id
	return r
}
func (r *Rect) GetId() any {
	return r.id
}

func (r *Rect) Rotate() *Rect {
	r.w, r.h = r.h, r.w
	return r
}

func (r *Rect) String() string {
	return fmt.Sprintf("[id: %v,w: %d, h: %d, ]", r.id, r.w, r.h)
}
