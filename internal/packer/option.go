package packer

import "os"

type options struct {
	UnpackFile  string // Unzip the file path
	InputDir    string // Enter a path
	OutputDir   string // Output path
	MaxW        int    // Maximum atlas width
	MaxH        int    // Maximum atlas height
	AutoSize    bool   // Automatically adjust atlas size
	Padding     int    // Padding
	Algorithm   int    // Packing algorithm
	Sort        bool   // Sorting
	AllowRotate bool   // Allow rotation
	Trim        bool   // Trim transparent pixels from the image
	Tolerance   uint8  // Tolerance for trimming transparency pixels 0-255
	SameDetect  bool   // Same detection
	PowerOfTwo  bool   // The atlas pixels are fixed to a power of 2
}

func NewOptions(optFns ...OptionFunc) *options {
	o := &options{
		UnpackFile:  "",
		InputDir:    "",
		OutputDir:   "",
		MaxW:        512,
		MaxH:        512,
		AutoSize:    false,
		Padding:     0,
		Algorithm:   AlgoBasic,
		Sort:        true,
		AllowRotate: false,
		Trim:        false,
		Tolerance:   0,
		SameDetect:  false,
		PowerOfTwo:  false,
	}
	for _, fn := range optFns {
		fn(o)
	}
	return o
}

type OptionFunc func(*options)

func WithUnpackFile(file string) OptionFunc {
	if file == "" {
		panic("unpack file path is empty")
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic("unpack file path is not exist")
	}
	return func(o *options) {
		o.UnpackFile = file
	}
}

func WithInputDir(dir string) OptionFunc {
	return func(o *options) {
		o.InputDir = dir
	}
}

func WithOutputDir(dir string) OptionFunc {
	return func(o *options) {
		o.OutputDir = dir
	}
}

func WithMaxSize(w, h int) OptionFunc {
	if w < 0 || h < 0 {
		panic("max size must be greater than 0")
	}
	return func(o *options) {
		o.MaxW = w
		o.MaxH = h
	}
}
func WithPadding(padding int) OptionFunc {
	if padding < 0 {
		panic("padding must be greater than 0")
	}
	return func(o *options) {
		o.Padding = padding
	}
}

func WithAlgorithm(algo int) OptionFunc {
	if algo < AlgoBasic || algo > AlgoMaxRectsBssf {
		algo = AlgoBasic
	}
	return func(o *options) {
		o.Algorithm = algo
	}
}
func WithSort(enable bool) OptionFunc {
	return func(o *options) {
		o.Sort = enable
	}
}
func WithAllowRotate(enable bool) OptionFunc {
	return func(o *options) {
		o.AllowRotate = enable
	}
}

func WithTrim(enable bool) OptionFunc {
	return func(o *options) {
		o.Trim = enable
	}
}

func WithTolerance(tolerance int) OptionFunc {
	if tolerance < 0 || tolerance > 255 {
		panic("tolerance must be between 0 and 255")
	}
	return func(o *options) {
		o.Tolerance = uint8(tolerance)
	}
}

func WithSameDetect(enable bool) OptionFunc {
	return func(o *options) {
		o.SameDetect = enable
	}
}
func WithPowerOfTwo(enable bool) OptionFunc {
	return func(o *options) {
		o.PowerOfTwo = enable
	}
}

type OptionBuilder struct {
	opt options
}

func NewOptionBuilder() *OptionBuilder {
	return &OptionBuilder{
		opt: options{
			UnpackFile:  "",
			InputDir:    "",
			OutputDir:   "",
			MaxW:        512,
			MaxH:        512,
			AutoSize:    false,
			Padding:     0,
			Algorithm:   AlgoBasic,
			Sort:        true,
			AllowRotate: false,
			Trim:        false,
			Tolerance:   0,
			SameDetect:  false,
			PowerOfTwo:  false,
		},
	}
}

func (b *OptionBuilder) UnpackFile(file string) *OptionBuilder {
	if file == "" {
		panic("unpack file path is empty")
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic("unpack file path is not exist")
	}
	b.opt.UnpackFile = file
	return b
}

func (b *OptionBuilder) InputDir(dir string) *OptionBuilder {
	b.opt.InputDir = dir
	return b
}
func (b *OptionBuilder) OutputDir(dir string) *OptionBuilder {
	b.opt.OutputDir = dir
	return b
}

func (b *OptionBuilder) MaxSize(w, h int) *OptionBuilder {
	if w < 0 || h < 0 {
		panic("max size must be greater than 0")
	}
	b.opt.MaxW = w
	b.opt.MaxH = h
	return b
}
func (b *OptionBuilder) Padding(padding int) *OptionBuilder {
	if padding < 0 {
		panic("padding must be greater than 0")
	}
	b.opt.Padding = padding
	return b
}
func (b *OptionBuilder) Algorithm(algo int) *OptionBuilder {
	if algo < AlgoBasic || algo > AlgoMaxRectsBssf {
		algo = AlgoBasic
	}
	b.opt.Algorithm = algo
	return b
}
func (b *OptionBuilder) Sort(enable bool) *OptionBuilder {
	b.opt.Sort = enable
	return b
}
func (b *OptionBuilder) AllowRotate(enable bool) *OptionBuilder {
	b.opt.AllowRotate = enable
	return b
}
func (b *OptionBuilder) Trim(enable bool) *OptionBuilder {
	b.opt.Trim = enable
	return b
}
func (b *OptionBuilder) Tolerance(tolerance int) *OptionBuilder {
	if tolerance < 0 || tolerance > 255 {
		panic("tolerance must be between 0 and 255")
	}
	b.opt.Tolerance = uint8(tolerance)
	return b
}
func (b *OptionBuilder) SameDetect(enable bool) *OptionBuilder {
	b.opt.SameDetect = enable
	return b
}
func (b *OptionBuilder) PowerOfTwo(enable bool) *OptionBuilder {
	b.opt.PowerOfTwo = enable
	return b
}
func (b *OptionBuilder) Build() *options {
	return &b.opt
}
