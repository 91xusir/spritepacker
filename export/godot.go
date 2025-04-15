package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/91xusir/spritepacker/pack"
	"text/template"
)

type gdAtlas struct {
	FileName string
	Width    int
	Height   int
}

type gdRect struct {
	Name   string
	Frame  pack.Rect
	Margin pack.Rect
	Last   bool
}

type godotTemplateData struct {
	Meta    pack.Meta `json:"meta"`
	Atlas   gdAtlas   `json:"Atlas"`
	GdRects []gdRect  `json:"rects"`
}

type GodotExporter struct {
	ext string
}

func (g *GodotExporter) Ext() string {
	return g.ext
}
func (g *GodotExporter) SetExt(ext string) {
	g.ext = ext
}

func (g *GodotExporter) Export(atlasInfo *pack.AtlasInfo) ([]byte, error) {
	if len(atlasInfo.Atlases) == 0 {
		return nil, fmt.Errorf("no Atlas found")
	}
	atlas := atlasInfo.Atlases[0]
	gdRects := make([]gdRect, len(atlas.Sprites))
	for i, sprite := range atlas.Sprites {
		margin := pack.Rect{}
		if sprite.Trimmed {
			margin.X = sprite.TrimmedRect.X
			margin.Y = sprite.Frame.Y - (sprite.TrimmedRect.Y)
			margin.W = sprite.SrcRect.W - sprite.Frame.W - margin.X
			margin.H = sprite.SrcRect.H - sprite.Frame.H - margin.Y
		}
		gdRects[i] = gdRect{
			Name:   sprite.FileName,
			Frame:  sprite.Frame,
			Margin: margin,
			Last:   i == len(atlas.Sprites)-1,
		}
	}
	data := godotTemplateData{
		Meta: atlasInfo.Meta,
		Atlas: gdAtlas{
			FileName: atlas.Name,
			Width:    atlas.Size.W,
			Height:   atlas.Size.H,
		},
		GdRects: gdRects,
	}

	tmpl, err := template.New(g.Ext()).Parse(godotTemplate)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.Bytes(), err
}

// Import TODO test
func (g *GodotExporter) Import(data []byte) (*pack.AtlasInfo, error) {
	var raw struct {
		Textures []struct {
			Image   string    `json:"image"`
			Size    pack.Size `json:"size"`
			Sprites []struct {
				Filename string    `json:"filename"`
				Region   pack.Rect `json:"region"`
				Margin   pack.Rect `json:"margin"`
			} `json:"sprites"`
		} `json:"textures"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}
	if len(raw.Textures) == 0 {
		return nil, fmt.Errorf("no texture found")
	}
	t := raw.Textures[0]
	sprites := make([]pack.Sprite, len(t.Sprites))
	for i, s := range t.Sprites {
		// 使用 margin 还原 trimmedRect 和 srcRect
		trimmed := s.Margin.X != 0 || s.Margin.Y != 0 || s.Margin.W != 0 || s.Margin.H != 0
		srcW := s.Region.W + s.Margin.X + s.Margin.W
		srcH := s.Region.H + s.Margin.Y + s.Margin.H
		trimmedRect := pack.NewRectByPosAndSize(s.Region.X-s.Margin.X, s.Region.Y-s.Margin.Y, srcW, srcH)
		sprites[i] = pack.Sprite{
			FileName:    s.Filename,
			Frame:       s.Region,
			SrcRect:     pack.Size{W: srcW, H: srcH},
			TrimmedRect: trimmedRect,
			Trimmed:     trimmed,
			Rotated:     false, // godot not support rotated
		}
	}
	return &pack.AtlasInfo{
		Meta: pack.Meta{
			Format: "tpsheet",
		},
		Atlases: []pack.Atlas{
			{
				Name:    t.Image,
				Size:    t.Size,
				Sprites: sprites,
			},
		},
	}, nil
}

const godotTemplate = `{
	"meta": {
		"repo": "{{.Meta.Repo}}",
		"format": "{{.Meta.Format}}",
		"version": "{{.Meta.Version}}",
		"timestamp": "{{.Meta.Timestamp}}"
	},
	"textures": [
		{
			"image": "{{.Atlas.FileName}}",
			"size": {
				"w": {{.Atlas.Width}},
				"h": {{.Atlas.Height}}
			},
			"sprites": [
				{{- range .GdRects}}
				{
					"filename": "{{.Name}}",
					"region": {
						"x": {{.Frame.X}},
						"y": {{.Frame.Y}},
						"w": {{.Frame.W}},
						"h": {{.Frame.H}}
					},
					"margin": {
						"x": {{.Margin.X}},
						"y": {{.Margin.Y}},
						"w": {{.Margin.W}},
						"h": {{.Margin.H}}
					}
				}{{if not .Last}},{{end}}{{end}}
			]
		}
	]
}`
