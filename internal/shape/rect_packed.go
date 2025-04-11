package shape

import "fmt"

// RectPacked represents a packed rectangle
type RectPacked struct {
	Rect
	X, Y      int // left-top point
	IsRotated bool
}

// NewRectPacked creates a new PackedRect with the given ID, x, y, width, height, and rotation flag.
func NewRectPacked(rect Rect, x, y int, isRotated bool) *RectPacked {
	if x < 0 || y < 0 {
		panic("packed rect x and y must be greater than 0")
	}
	return &RectPacked{
		Rect:      rect,
		X:         x,
		Y:         y,
		IsRotated: isRotated,
	}
}

func (r *RectPacked) String() string {
	return fmt.Sprintf("[id: %v, x: %d, y: %d, w: %d, h: %d, isRotated: %v]", r.id, r.X, r.Y, r.w, r.h, r.IsRotated)
}
