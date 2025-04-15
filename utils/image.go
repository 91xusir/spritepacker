package utils

import (
	"errors"
	"fmt"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

//----------------image--------------------

type ImageFormat int

const (
	JPEG ImageFormat = iota
	PNG
	TIFF
	BMP
	WEBP
)

// Compression level
type clv int

const (
	NoCompression      clv = 0 // No compression
	BestSpeed          clv = 1 // Best speed
	DefaultCompression clv = 2 // Default compression
	BestCompression    clv = 3 // Best compression
)

type encodeImgOpt struct {
	jpegQuality         int
	pngCompressionLevel png.CompressionLevel
	ttfCompressionLevel tiff.CompressionType
	ttfPredictor        bool
}
type SetClv func(*encodeImgOpt)

// WithCLV sets the compression level for the image.
func WithCLV(clv clv) SetClv {
	return func(opts *encodeImgOpt) {
		switch clv {
		case NoCompression:
			opts.jpegQuality = 100
			opts.pngCompressionLevel = png.NoCompression
			opts.ttfCompressionLevel = tiff.Uncompressed
		case BestSpeed:
			opts.jpegQuality = 85
			opts.pngCompressionLevel = png.BestSpeed
			opts.ttfCompressionLevel = tiff.LZW
			opts.ttfPredictor = false
		case DefaultCompression:
			opts.jpegQuality = jpeg.DefaultQuality
			opts.pngCompressionLevel = png.DefaultCompression
			opts.ttfCompressionLevel = tiff.Deflate
			opts.ttfPredictor = true
		case BestCompression:
			opts.jpegQuality = 60
			opts.pngCompressionLevel = png.BestCompression
			opts.ttfCompressionLevel = tiff.Deflate
			opts.ttfPredictor = true
		default:
			opts.jpegQuality = 75
			opts.pngCompressionLevel = png.DefaultCompression
			opts.ttfCompressionLevel = tiff.Deflate
			opts.ttfPredictor = true
		}
	}
}

func SaveImgByExt(outputPath string, img image.Image, compressionLevel ...SetClv) error {
	outputPath = strings.ToLower(outputPath)
	ext := filepath.Ext(outputPath)
	if ext == "" {
		ext = ".png"
	}
	var format ImageFormat
	switch ext {
	case ".jpg", ".jpeg":
		format = JPEG
	case ".png":
		format = PNG
	case ".tiff":
		format = TIFF
	case ".bmp":
		format = BMP
	default:
		return errors.New("unsupported image format")
	}
	file, err := SafeCreate(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return EncImg(file, img, format, compressionLevel...)
}

func SaveImgAutoExt(pathNoneExt string, img image.Image, format ImageFormat, compressionLevel ...SetClv) error {
	opts := &encodeImgOpt{
		jpegQuality:         jpeg.DefaultQuality,
		pngCompressionLevel: png.DefaultCompression,
		ttfCompressionLevel: tiff.Deflate,
		ttfPredictor:        true,
	}
	for _, clv := range compressionLevel {
		clv(opts)
	}
	var ext string
	switch format {
	case JPEG:
		ext = ".jpeg"
	case PNG:
		ext = ".png"
	case TIFF:
		ext = ".tiff"
	case BMP:
		ext = ".bmp"
	default:
		return errors.New("unsupported image format")
	}
	fullPath := pathNoneExt + ext
	file, err := SafeCreate(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return EncImg(file, img, format, compressionLevel...)
}

// EncImg encodes the image to the specified format and writes it to the writer.
// Parameters:
//   - w: the writer to write the encoded image
//   - img: the image to encode
//   - format: the image format to encode
//   - compressionLevel: the compression level to use, default is NoCompression
//
// Returns:
//   - error: an error if the format is unsupported or if encoding fails
//
// Supported formats:
//   - JPEG
//   - PNG
//   - TIFF
//   - BMP
func EncImg(w io.Writer, img image.Image, format ImageFormat, compressionLevel ...SetClv) error {
	opts := &encodeImgOpt{
		jpegQuality:         jpeg.DefaultQuality,
		pngCompressionLevel: png.DefaultCompression,
		ttfCompressionLevel: tiff.Deflate,
		ttfPredictor:        true,
	}
	for _, clv := range compressionLevel {
		clv(opts)
	}

	switch format {
	case JPEG:
		if nrgba, ok := img.(*image.NRGBA); ok && nrgba.Opaque() {
			rgba := &image.RGBA{
				Pix:    nrgba.Pix,
				Stride: nrgba.Stride,
				Rect:   nrgba.Rect,
			}
			return jpeg.Encode(w, rgba, &jpeg.Options{Quality: opts.jpegQuality})
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: opts.jpegQuality})
	case PNG:
		encoder := png.Encoder{CompressionLevel: opts.pngCompressionLevel}
		return encoder.Encode(w, img)
	case TIFF:
		return tiff.Encode(w, img, &tiff.Options{
			Compression: opts.ttfCompressionLevel,
			Predictor:   opts.ttfPredictor,
		})
	case BMP:
		return bmp.Encode(w, img)
	default:
		return errors.New("unsupported image format")
	}
}

func LoadImg(pathName string) (image.Image, error) {
	file, err := os.Open(pathName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return DecImg(file)
}

// DecImg
// Supported formats:
//   - JPEG
//   - PNG
//   - TIFF
//   - BMP
//   - WEBP
func DecImg(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// GetOpaqueBounds returns the bounds of the opaque area of the image.
// The image is assumed to be in RGBA format.
func GetOpaqueBounds(img image.Image, tolerance uint8) image.Rectangle {
	bounds := img.Bounds()
	if bounds.Empty() {
		return bounds
	}

	// Initialize to opposite extremes
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y
	found := false

	// Helper function to update bounds
	updateBounds := func(x, y int) {
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		found = true
	}

	switch src := img.(type) {
	case *image.RGBA, *image.NRGBA:
		// Unified handling for RGBA and NRGBA
		var pix []uint8
		var stride int
		switch t := src.(type) {
		case *image.RGBA:
			pix = t.Pix
			stride = t.Stride
		case *image.NRGBA:
			pix = t.Pix
			stride = t.Stride
		}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			i := (y-bounds.Min.Y)*stride + (bounds.Min.X-bounds.Min.X)*4
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if pix[i+3] > tolerance {
					updateBounds(x, y)
				}
				i += 4
			}
		}

	default:
		// Generic case for other image types
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				_, _, _, a := img.At(x, y).RGBA()
				if uint8(a>>8) > tolerance { // Convert 16-bit alpha to 8-bit
					updateBounds(x, y)
				}
			}
		}
	}

	if !found {
		return bounds
	}
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

// Rotate90 rotates the image 90 degrees counter-clockwise and returns the transformed image.
func Rotate90(img image.Image) *image.NRGBA {
	src := newScanner(img)
	dstW := src.h
	dstH := src.w
	rowSize := dstW * 4
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))
	Parallel(0, dstH, func(ys <-chan int) {
		for dstY := range ys {
			i := dstY * dst.Stride
			srcX := dstH - dstY - 1
			src.scan(srcX, 0, srcX+1, src.h, dst.Pix[i:i+rowSize])
		}
	})
	return dst
}

// Rotate180 rotates the image 180 degrees counter-clockwise and returns the transformed image.
func Rotate180(img image.Image) *image.NRGBA {
	src := newScanner(img)
	dstW := src.w
	dstH := src.h
	rowSize := dstW * 4
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))
	Parallel(0, dstH, func(ys <-chan int) {
		for dstY := range ys {
			i := dstY * dst.Stride
			srcY := dstH - dstY - 1
			src.scan(0, srcY, src.w, srcY+1, dst.Pix[i:i+rowSize])
			reverse(dst.Pix[i : i+rowSize])
		}
	})
	return dst
}

// Rotate270 rotates the image 270 degrees counter-clockwise and returns the transformed image.
func Rotate270(img image.Image) *image.NRGBA {
	src := newScanner(img)
	dstW := src.h
	dstH := src.w
	rowSize := dstW * 4
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))
	Parallel(0, dstH, func(ys <-chan int) {
		for dstY := range ys {
			i := dstY * dst.Stride
			srcX := dstY
			src.scan(srcX, 0, srcX+1, src.h, dst.Pix[i:i+rowSize])
			reverse(dst.Pix[i : i+rowSize])
		}
	})
	return dst
}

func isImageDifferent(img1Path, img2Path string) bool {
	file1, err := os.Open(img1Path)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", img1Path, err)
		return true
	}
	defer file1.Close()

	file2, err := os.Open(img2Path)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", img2Path, err)
		return true
	}
	defer file2.Close()

	img1, err := DecImg(file1)
	if err != nil {
		fmt.Printf("Failed to decode %s: %v\n", img1Path, err)
		return true
	}

	img2, err := DecImg(file2)
	if err != nil {
		fmt.Printf("Failed to decode %s: %v\n", img2Path, err)
		return true
	}

	if img1.Bounds().Size() != img2.Bounds().Size() {
		return true
	}

	bounds := img1.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img1.At(x, y) != img2.At(x, y) {
				return true
			}
		}
	}

	return false
}

func CompareImgFormFolders(inputDir, outputDir string) []string {
	var differentFiles []string
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Println("Error reading input dir:", err)
		return nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		inputPath := filepath.Join(inputDir, entry.Name())
		outputPath := filepath.Join(outputDir, entry.Name())
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			fmt.Printf("Missing file: %s\n", outputPath)
			continue
		}
		if isImageDifferent(inputPath, outputPath) {
			differentFiles = append(differentFiles, entry.Name())
		}
	}
	return differentFiles
}

type scanner struct {
	image   image.Image
	w, h    int
	palette []color.NRGBA
}

func newScanner(img image.Image) *scanner {
	s := &scanner{
		image: img,
		w:     img.Bounds().Dx(),
		h:     img.Bounds().Dy(),
	}
	if img, ok := img.(*image.Paletted); ok {
		s.palette = make([]color.NRGBA, len(img.Palette))
		for i := 0; i < len(img.Palette); i++ {
			s.palette[i] = color.NRGBAModel.Convert(img.Palette[i]).(color.NRGBA)
		}
	}
	return s
}

func (s *scanner) scan(x1, y1, x2, y2 int, dst []uint8) {
	switch img := s.image.(type) {
	case *image.NRGBA:
		size := (x2 - x1) * 4
		j := 0
		i := y1*img.Stride + x1*4
		if size == 4 {
			for y := y1; y < y2; y++ {
				d := dst[j : j+4 : j+4]
				s := img.Pix[i : i+4 : i+4]
				d[0] = s[0]
				d[1] = s[1]
				d[2] = s[2]
				d[3] = s[3]
				j += size
				i += img.Stride
			}
		} else {
			for y := y1; y < y2; y++ {
				copy(dst[j:j+size], img.Pix[i:i+size])
				j += size
				i += img.Stride
			}
		}

	case *image.NRGBA64:
		j := 0
		for y := y1; y < y2; y++ {
			i := y*img.Stride + x1*8
			for x := x1; x < x2; x++ {
				s := img.Pix[i : i+8 : i+8]
				d := dst[j : j+4 : j+4]
				d[0] = s[0]
				d[1] = s[2]
				d[2] = s[4]
				d[3] = s[6]
				j += 4
				i += 8
			}
		}

	case *image.RGBA:
		j := 0
		for y := y1; y < y2; y++ {
			i := y*img.Stride + x1*4
			for x := x1; x < x2; x++ {
				d := dst[j : j+4 : j+4]
				a := img.Pix[i+3]
				switch a {
				case 0:
					d[0] = 0
					d[1] = 0
					d[2] = 0
					d[3] = a
				case 0xff:
					s := img.Pix[i : i+4 : i+4]
					d[0] = s[0]
					d[1] = s[1]
					d[2] = s[2]
					d[3] = a
				default:
					s := img.Pix[i : i+4 : i+4]
					r16 := uint16(s[0])
					g16 := uint16(s[1])
					b16 := uint16(s[2])
					a16 := uint16(a)
					d[0] = uint8(r16 * 0xff / a16)
					d[1] = uint8(g16 * 0xff / a16)
					d[2] = uint8(b16 * 0xff / a16)
					d[3] = a
				}
				j += 4
				i += 4
			}
		}

	case *image.RGBA64:
		j := 0
		for y := y1; y < y2; y++ {
			i := y*img.Stride + x1*8
			for x := x1; x < x2; x++ {
				s := img.Pix[i : i+8 : i+8]
				d := dst[j : j+4 : j+4]
				a := s[6]
				switch a {
				case 0:
					d[0] = 0
					d[1] = 0
					d[2] = 0
				case 0xff:
					d[0] = s[0]
					d[1] = s[2]
					d[2] = s[4]
				default:
					r32 := uint32(s[0])<<8 | uint32(s[1])
					g32 := uint32(s[2])<<8 | uint32(s[3])
					b32 := uint32(s[4])<<8 | uint32(s[5])
					a32 := uint32(s[6])<<8 | uint32(s[7])
					d[0] = uint8((r32 * 0xffff / a32) >> 8)
					d[1] = uint8((g32 * 0xffff / a32) >> 8)
					d[2] = uint8((b32 * 0xffff / a32) >> 8)
				}
				d[3] = a
				j += 4
				i += 8
			}
		}

	case *image.Gray:
		j := 0
		for y := y1; y < y2; y++ {
			i := y*img.Stride + x1
			for x := x1; x < x2; x++ {
				c := img.Pix[i]
				d := dst[j : j+4 : j+4]
				d[0] = c
				d[1] = c
				d[2] = c
				d[3] = 0xff
				j += 4
				i++
			}
		}

	case *image.Gray16:
		j := 0
		for y := y1; y < y2; y++ {
			i := y*img.Stride + x1*2
			for x := x1; x < x2; x++ {
				c := img.Pix[i]
				d := dst[j : j+4 : j+4]
				d[0] = c
				d[1] = c
				d[2] = c
				d[3] = 0xff
				j += 4
				i += 2
			}
		}

	case *image.YCbCr:
		j := 0
		x1 += img.Rect.Min.X
		x2 += img.Rect.Min.X
		y1 += img.Rect.Min.Y
		y2 += img.Rect.Min.Y

		hy := img.Rect.Min.Y / 2
		hx := img.Rect.Min.X / 2
		for y := y1; y < y2; y++ {
			iy := (y-img.Rect.Min.Y)*img.YStride + (x1 - img.Rect.Min.X)

			var yBase int
			switch img.SubsampleRatio {
			case image.YCbCrSubsampleRatio444, image.YCbCrSubsampleRatio422:
				yBase = (y - img.Rect.Min.Y) * img.CStride
			case image.YCbCrSubsampleRatio420, image.YCbCrSubsampleRatio440:
				yBase = (y/2 - hy) * img.CStride
			}

			for x := x1; x < x2; x++ {
				var ic int
				switch img.SubsampleRatio {
				case image.YCbCrSubsampleRatio444, image.YCbCrSubsampleRatio440:
					ic = yBase + (x - img.Rect.Min.X)
				case image.YCbCrSubsampleRatio422, image.YCbCrSubsampleRatio420:
					ic = yBase + (x/2 - hx)
				default:
					ic = img.COffset(x, y)
				}

				yy1 := int32(img.Y[iy]) * 0x10101
				cb1 := int32(img.Cb[ic]) - 128
				cr1 := int32(img.Cr[ic]) - 128

				r := yy1 + 91881*cr1
				if uint32(r)&0xff000000 == 0 {
					r >>= 16
				} else {
					r = ^(r >> 31)
				}

				g := yy1 - 22554*cb1 - 46802*cr1
				if uint32(g)&0xff000000 == 0 {
					g >>= 16
				} else {
					g = ^(g >> 31)
				}

				b := yy1 + 116130*cb1
				if uint32(b)&0xff000000 == 0 {
					b >>= 16
				} else {
					b = ^(b >> 31)
				}

				d := dst[j : j+4 : j+4]
				d[0] = uint8(r)
				d[1] = uint8(g)
				d[2] = uint8(b)
				d[3] = 0xff

				iy++
				j += 4
			}
		}

	case *image.Paletted:
		j := 0
		for y := y1; y < y2; y++ {
			i := y*img.Stride + x1
			for x := x1; x < x2; x++ {
				c := s.palette[img.Pix[i]]
				d := dst[j : j+4 : j+4]
				d[0] = c.R
				d[1] = c.G
				d[2] = c.B
				d[3] = c.A
				j += 4
				i++
			}
		}

	default:
		j := 0
		b := s.image.Bounds()
		x1 += b.Min.X
		x2 += b.Min.X
		y1 += b.Min.Y
		y2 += b.Min.Y
		for y := y1; y < y2; y++ {
			for x := x1; x < x2; x++ {
				r16, g16, b16, a16 := s.image.At(x, y).RGBA()
				d := dst[j : j+4 : j+4]
				switch a16 {
				case 0xffff:
					d[0] = uint8(r16 >> 8)
					d[1] = uint8(g16 >> 8)
					d[2] = uint8(b16 >> 8)
					d[3] = 0xff
				case 0:
					d[0] = 0
					d[1] = 0
					d[2] = 0
					d[3] = 0
				default:
					d[0] = uint8(((r16 * 0xffff) / a16) >> 8)
					d[1] = uint8(((g16 * 0xffff) / a16) >> 8)
					d[2] = uint8(((b16 * 0xffff) / a16) >> 8)
					d[3] = uint8(a16 >> 8)
				}
				j += 4
			}
		}
	}
}

func reverse(pix []uint8) {
	if len(pix) <= 4 {
		return
	}
	i := 0
	j := len(pix) - 4
	for i < j {
		pi := pix[i : i+4 : i+4]
		pj := pix[j : j+4 : j+4]
		pi[0], pj[0] = pj[0], pi[0]
		pi[1], pj[1] = pj[1], pi[1]
		pi[2], pj[2] = pj[2], pi[2]
		pi[3], pj[3] = pj[3], pi[3]
		i += 4
		j -= 4
	}
}

// Parallel  processes the data in separate goroutines.
func Parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}
	process := runtime.GOMAXPROCS(0)
	if process > count {
		process = count
	}

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)
	var wg sync.WaitGroup
	for i := 0; i < process; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}
