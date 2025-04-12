package pack

import (
	"fmt"
	"os"
)

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
	//----validate----
	err error
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
	if b.err != nil {
		return b
	}
	if dir == "" {
		b.err = fmt.Errorf("input dir is empty")
		return b
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		b.err = fmt.Errorf("input dir does not exist: %w", err)
		return b
	}
	b.inputDir = dir
	return b
}

func (b *Options) OutputDir(dir string) *Options {
	if b.err != nil {
		return b
	}
	if dir == "" {
		b.err = fmt.Errorf("output dir is empty")
		return b
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		b.err = fmt.Errorf("output dir does not exist: %w", err)
		return b
	}
	b.outputDir = dir
	return b
}

func (b *Options) MaxSize(w, h int) *Options {
	if b.err != nil {
		return b
	}
	if w <= 0 || h <= 0 {
		b.err = fmt.Errorf("max size must be greater than 0")
		return b
	}
	b.maxW = w
	b.maxH = h
	return b
}

func (b *Options) Padding(padding int) *Options {
	if b.err != nil {
		return b
	}
	if padding < 0 {
		b.err = fmt.Errorf("padding must be >= 0")
		return b
	}
	b.padding = padding
	return b
}

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
func (b *Options) Sort(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.sort = enable
	return b
}
func (b *Options) AllowRotate(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.allowRotate = enable
	return b
}
func (b *Options) Trim(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.trim = enable
	return b
}
func (b *Options) AutoSize(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.autoSize = enable
	return b
}
func (b *Options) Tolerance(tolerance int) *Options {
	if b.err != nil {
		return b
	}
	if tolerance < 0 || tolerance > 255 {
		b.err = fmt.Errorf("tolerance must be in the range 0-255")
		return b
	}
	b.tolerance = uint8(tolerance)
	return b
}
func (b *Options) SameDetect(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.sameDetect = enable
	return b
}
func (b *Options) PowerOfTwo(enable bool) *Options {
	if b.err != nil {
		return b
	}
	b.powerOfTwo = enable
	return b
}

func (b *Options) Validate() (*Options, error) {
	return b, b.err
}
