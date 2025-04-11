package texture

import "image"


// GetOpaqueBounds returns the bounds of the opaque area of the image.
// The image is assumed to be in RGBA format.
func GetOpaqueBounds(img image.Image, alphaThreshold uint32) image.Rectangle {
	bounds := img.Bounds()
	if bounds.Empty() {
		return bounds
	}
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y
	found := false
	switch src := img.(type) {
	case *image.RGBA:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			i := src.PixOffset(bounds.Min.X, y)
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if src.Pix[i+3] > uint8(alphaThreshold) { 
					found = true
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
				}
				i += 4
			}
		}
	case *image.NRGBA:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			i := src.PixOffset(bounds.Min.X, y)
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if src.Pix[i+3] > uint8(alphaThreshold) {
					found = true
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
				}
				i += 4
			}
		}
	default:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				_, _, _, a := img.At(x, y).RGBA()
				a8 := a >> 8
				if a8 > alphaThreshold {
					found = true
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
				}
			}
		}
	}
	if !found {
		return bounds 
	}
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

