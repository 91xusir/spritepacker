package pack

import "os"

type Options struct {
	//----path----
	inputDir  string // enter a path
	outputDir string // output path
	//----rect----
	maxW        int       // maximum atlas width
	maxH        int       // maximum atlas height
	autoSize    bool      // automatically adjust atlas size
	padding     int       // padding
	algorithm   Algorithm // packing algorithm
	heuristic   Heuristic // heuristic it is valid only when the algorithm is AlgoMaxRects
	allowRotate bool      // allow rotation
	//----image----
	sort       bool  // sorting by file name
	trim       bool  // trim transparent pixels from the image
	tolerance  uint8 // tolerance for trimming transparency pixels 0-255
	sameDetect bool  // same detection
	powerOfTwo bool  // the atlas pixels are fixed to a power of 2
}

func NewOptions() *Options {
	return &Options{
		inputDir:    "",
		outputDir:   "",
		maxW:        512,
		maxH:        512,
		autoSize:    false,
		padding:     0,
		algorithm:   AlgoBasic,
		heuristic:   BestShortSideFit,
		sort:        true,
		allowRotate: false,
		trim:        false,
		tolerance:   0,
		sameDetect:  false,
		powerOfTwo:  false,
	}
}

func (b *Options) InputDir(dir string) *Options {
	if dir == "" {
		panic("input dir is empty")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		panic("input dir does not exist")
	}
	b.inputDir = dir
	return b
}
func (b *Options) OutputDir(dir string) *Options {
	if dir == "" {
		panic("output dir is empty")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		panic("output dir does not exist")
	}
	b.outputDir = dir
	return b
}

func (b *Options) MaxSize(w, h int) *Options {
	if w <= 0 || h <= 0 {
		panic("max size must be greater than 0")
	}
	b.maxW = w
	b.maxH = h
	return b
}
func (b *Options) Padding(padding int) *Options {
	if padding < 0 {
		panic("padding must be >= 0")
	}
	b.padding = padding
	return b
}
func (b *Options) Algorithm(algo Algorithm) *Options {
	if algo < AlgoBasic || algo >= MaxAlgosIndex {
		algo = AlgoBasic
	}
	b.algorithm = algo
	return b
}
func (b *Options) Heuristic(heuristic Heuristic) *Options {
	if heuristic < BestShortSideFit || heuristic >= MaxHeuristicsIndex {
		heuristic = BestShortSideFit
	}
	b.heuristic = heuristic
	return b
}
func (b *Options) Sort(enable bool) *Options {
	b.sort = enable
	return b
}
func (b *Options) AllowRotate(enable bool) *Options {
	b.allowRotate = enable
	return b
}
func (b *Options) Trim(enable bool) *Options {
	b.trim = enable
	return b
}
func (b *Options) AutoSize(enable bool) *Options {
	b.autoSize = enable
	return b
}
func (b *Options) Tolerance(tolerance int) *Options {
	if tolerance < 0 || tolerance > 255 {
		panic("tolerance must be between 0 and 255")
	}
	b.tolerance = uint8(tolerance)
	return b
}
func (b *Options) SameDetect(enable bool) *Options {
	b.sameDetect = enable
	return b
}
func (b *Options) PowerOfTwo(enable bool) *Options {
	b.powerOfTwo = enable
	return b
}
