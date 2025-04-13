package pack

import (
	"errors"
)

type Options struct {
	//----rect----
	maxW        int       // maximum atlas width
	maxH        int       // maximum atlas height
	autoSize    bool      // automatically adjust atlas size
	padding     int       // padding
	algorithm   Algorithm // packing algorithm
	heuristic   Heuristic // heuristic it is valid only when the algorithm is AlgoMaxRects
	allowRotate bool      // allow rotation

	//----atlas----
	name               string // atlas name
	sort               bool   // sorting by file name
	trim               bool   // trim transparent pixels from the image
	tolerance          uint8  // tolerance for trimming transparency pixels 0-255
	sameDetect         bool   // same detection
	powerOfTwo         bool   // the atlas pixels are fixed to a power of 2
	preMultipliedAlpha bool   // premultiplied alpha
	//----validate----
	err error
}

func NewOptions() *Options {
	return &Options{
		maxW:               512,
		maxH:               512,
		name:               "atlas",
		autoSize:           false,
		padding:            0,
		algorithm:          AlgoBasic,
		heuristic:          BestShortSideFit,
		sort:               true,
		allowRotate:        false,
		trim:               false,
		tolerance:          0,
		sameDetect:         false,
		preMultipliedAlpha: false,
		powerOfTwo:         false,
	}
}

// MaxSize sets the maximum size of the atlas.
// If the width or height is less than or equal to 0, it will be set to 512.
func (b *Options) MaxSize(w, h int) *Options {
	if b.err != nil {
		return b
	}
	if w <= 0 || h <= 0 {
		b.err = errors.New("max size must be greater than 0")
		return b
	}
	b.maxW = w
	b.maxH = h
	return b
}

// Name sets the name of the atlas.
// If the name is empty, the default name is "atlas".
func (b *Options) Name(name string) *Options {
	if b.err != nil {
		return b
	}
	if name == "" {
		name = "atlas"
	}
	b.name = name
	return b
}

// Padding sets the padding of the atlas.
// If the padding is less than 0, it will be set to 0.
func (b *Options) Padding(padding int) *Options {
	if b.err != nil {
		return b
	}
	if padding < 0 {
		b.err = errors.New("padding must be >= 0")
		return b
	}
	b.padding = padding
	return b
}

// Algorithm sets the packing algorithm of the atlas.
// If the algorithm is not valid, it will be set to AlgoBasic.
func (b *Options) Algorithm(algo Algorithm) *Options {
	if b.err != nil {
		return b
	}
	if algo < AlgoBasic || algo >= MaxAlgoIndex {
		algo = AlgoBasic
	}
	b.algorithm = algo
	return b
}

// Heuristic sets the heuristic of the atlas.
// If the heuristic is not valid, it will be set to BestShortSideFit.
// It is valid only when the algorithm is AlgoMaxRects.
func (b *Options) Heuristic(heuristic Heuristic) *Options {
	if b.err != nil {
		return b
	}
	if heuristic < BestShortSideFit || heuristic >= MaxHeuristicsIndex {
		heuristic = BestShortSideFit
	}
	b.heuristic = heuristic
	return b
}

// Sort sets the sorting of the atlas.
// default method is sort by area.
func (b *Options) Sort(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.sort = enable
	return b
}

// AllowRotate allow sprites to Rotate to optimizeAtlasSize
func (b *Options) AllowRotate(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.allowRotate = enable
	return b
}

// Trim transparent pixels will be trimmed from the image.
func (b *Options) Trim(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.trim = enable
	return b
}

// Tolerance sets the tolerance of trimming transparency pixels.
func (b *Options) Tolerance(tolerance int) *Options {
	if b.err != nil {
		return b
	}
	if tolerance < 0 || tolerance > 255 {
		b.err = errors.New("tolerance must be in the range 0-255")
		return b
	}
	b.tolerance = uint8(tolerance)
	return b
}

// AutoSize atlas size will be automatically adjusted.
func (b *Options) AutoSize(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.autoSize = enable
	return b
}

// SameDetect sets the same detection of the atlas.
func (b *Options) SameDetect(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.sameDetect = enable
	return b
}

// PowerOfTwo sets the power of two of the atlas.
// The atlas pixels are fixed to a power of 2.
func (b *Options) PowerOfTwo(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.powerOfTwo = enable
	return b
}

// PmAlpha sets the premultiplied alpha.
func (b *Options) PmAlpha(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.preMultipliedAlpha = enable
	return b
}

// Validate validates the options.
// If the options are invalid, it will return an error.
func (b *Options) Validate() (*Options, error) {
	return b, b.err
}
