package pack

import (
	"fmt"
	"github.com/91xusir/spritepacker/export"
	"github.com/91xusir/spritepacker/model"
	"github.com/91xusir/spritepacker/utils"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"
)

type unpackedOpts struct {
	atlasImgPath string
	outputPath   string
}
type UnpackOpts func(*unpackedOpts)

func WithImgInput(atlasImgPath string) UnpackOpts {
	if atlasImgPath == "" {
		return func(opts *unpackedOpts) {
			return
		}
	}
	return func(opts *unpackedOpts) {
		opts.atlasImgPath = atlasImgPath
	}
}
func WithOutput(outputPath string) UnpackOpts {
	if outputPath == "" {
		return func(opts *unpackedOpts) {
			return
		}
	}
	return func(opts *unpackedOpts) {
		opts.outputPath = outputPath
	}
}

func UnpackSprites(infoPath string, fn ...UnpackOpts) error {
	opts := &unpackedOpts{
		atlasImgPath: filepath.Dir(infoPath),
		outputPath:   filepath.Dir(infoPath),
	}
	for _, f := range fn {
		f(opts)
	}
	// make sure outputPath exists
	if err := os.MkdirAll(opts.outputPath, os.ModePerm); err != nil {
		return err
	}
	exporter := export.NewExportManager().Init()
	atlasInfo, err := exporter.Import(infoPath)
	if err != nil {
		return err
	}
	err = UnpackAtlas(atlasInfo, *opts)
	if err != nil {
		return err
	}
	return nil
}

func UnpackAtlas(atlasInfo *model.AtlasInfo, opts unpackedOpts) error {
	baseNames := make([]string, len(atlasInfo.Atlases))
	for i := range atlasInfo.Atlases {
		baseNames[i] = strings.TrimSuffix(filepath.Base(atlasInfo.Atlases[i].Name), filepath.Ext(atlasInfo.Atlases[i].Name))
	}
	extStr := []string{".png", ".jpg", ".jpeg", ".bmp", ".tiff", ".webp"}
	for i, baseName := range baseNames {
		var imgFilePath string
		found := false
		for _, ext := range extStr {
			tryPath := filepath.Join(opts.atlasImgPath, baseName+ext)
			if _, err := os.Stat(tryPath); err == nil {
				imgFilePath = tryPath
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("image file for atlas %s not found", baseName)
		}
		atlasImg, err := utils.LoadImg(imgFilePath)
		if err != nil {
			return fmt.Errorf("failed to load image %s: %v", imgFilePath, err)
		}
		for j := range atlasInfo.Atlases[i].Sprites {
			sprite := atlasInfo.Atlases[i].Sprites[j]
			outputPath := filepath.Join(opts.outputPath, filepath.Base(sprite.FileName))
			subImg := image.NewNRGBA(image.Rect(0, 0, sprite.Frame.W, sprite.Frame.H))
			srcLeftTopPoint := image.Point{
				X: sprite.Frame.X,
				Y: sprite.Frame.Y,
			}
			draw.Draw(subImg, subImg.Bounds(), atlasImg, srcLeftTopPoint, draw.Src)
			// if rotated
			if sprite.Rotated {
				subImg = utils.Rotate90(subImg)
			}
			// if trimmed
			if sprite.Trimmed {
				img := image.NewNRGBA(image.Rect(0, 0, sprite.SrcRect.W, sprite.SrcRect.H))
				destRect := image.Rect(
					sprite.TrimmedRect.X,
					sprite.TrimmedRect.Y,
					sprite.TrimmedRect.X+subImg.Bounds().Dx(),
					sprite.TrimmedRect.Y+subImg.Bounds().Dy(),
				)
				draw.Draw(img, destRect, subImg, image.Point{}, draw.Src)
				subImg = img
			}
			err := utils.SaveImgByExt(outputPath, subImg)
			if err != nil {
				return fmt.Errorf("failed to save image %s: %v", outputPath, err)
			}
		}
	}
	return nil
}
