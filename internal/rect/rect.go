package rect

// Rect represents a rectangular structure
// with an ID, width, and height.
type Rect struct {
	Name string
	W    int
	H    int
}

// NewRect creates a new Rect with the given ID, width, and height.
// It returns a pointer to the newly created Rect.
func NewRect(name string, w int, h int) *Rect {
	return &Rect{
		Name: name,
		W:    w,
		H:    h,
	}
}

// Clone creates a new Rect with the same ID, width, and height as the given Rect.
// It returns a pointer to the newly created Rect.
func Clone(r *Rect) *Rect {
	return NewRect(r.Name, r.W, r.H)
}

// CloneRects creates a new slice of Rects with the same ID, width, and height as the given slice of Rects.
// It returns a new slice of Rects.
func CloneRects(rs []*Rect) []*Rect {
	newRects := make([]*Rect, len(rs))
	for i, item := range rs {
		newRects[i] = Clone(item)
	}
	return newRects
}
