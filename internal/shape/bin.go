package shape

// Bin represents a Bin with a width, height, and a list of rectangles to be packed.
// It also has a flag to indicate whether rectangle rotation is allowed.
type Bin struct {
	Rect

	// Rectangle list
	ReqPackRectList []*Rect
}

// NewBin creates a new bin with the given width, height, and list of rectangles to be packed.
func NewBin(w, h int, req []*Rect) *Bin {
	return NewBinById(w, h, req, nil)
}

// NewBinById creates a new bin with the given width, height, list of rectangles to be packed, and ID.
func NewBinById(w, h int, req []*Rect, id any) *Bin {
	bin := Bin{
		ReqPackRectList: req,
	}
	bin.SetH(h).SetW(w).SetId(id)
	return &bin
}
