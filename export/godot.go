package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/91xusir/spritepacker/pack"
	"text/template"
)

// 模板数据结构
type godotTemplateData struct {
	Config godotConfig `json:"config"`
	Rects  []godotRect `json:"rects"`
	Mata   pack.Meta   `json:"mata"`
}

type godotConfig struct {
	ImageFile   string
	ImageWidth  int
	ImageHeight int
}

type godotRect struct {
	Name   string
	Frame  pack.Rect
	Margin pack.Rect
	Last   bool
}

type GodotExporter struct{}

func (g *GodotExporter) Export(atlas *pack.SpriteAtlas) ([]byte, error) {
	if len(atlas.Atlases) == 0 {
		return nil, fmt.Errorf("no atlas found")
	}
	a := atlas.Atlases[0]
	rects := make([]godotRect, len(a.Sprites))
	for i, sprite := range a.Sprites {
		margin := pack.Rect{}
		if sprite.Trimmed {
			margin.X = sprite.Frame.X - (sprite.TrimmedRect.X)
			margin.Y = sprite.Frame.Y - (sprite.TrimmedRect.Y)
			margin.W = sprite.SrcRect.W - sprite.Frame.W - margin.X
			margin.H = sprite.SrcRect.H - sprite.Frame.H - margin.Y
		}
		rects[i] = godotRect{
			Name:   sprite.FileName,
			Frame:  sprite.Frame,
			Margin: margin,
			Last:   i == len(a.Sprites)-1,
		}
	}
	data := godotTemplateData{
		Config: godotConfig{
			ImageFile:   a.Name,
			ImageWidth:  a.Size.W,
			ImageHeight: a.Size.H,
		},
		Rects: rects,
		Mata:  atlas.Meta,
	}

	tmpl, err := template.New("tpsheet").Parse(godotTemplate)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.Bytes(), err
}

func (g *GodotExporter) Import(data []byte) (*pack.SpriteAtlas, error) {
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
		trimmedRect := pack.Rect{
			X: s.Region.X - s.Margin.X,
			Y: s.Region.Y - s.Margin.Y,
			W: srcW,
			H: srcH,
		}
		sprites[i] = pack.Sprite{
			FileName:    s.Filename,
			Frame:       s.Region,
			SrcRect:     Size{W: srcW, H: srcH},
			TrimmedRect: trimmedRect,
			Trimmed:     trimmed,
			Rotated:     false, // 可扩展支持
		}
	}
	return &pack.SpriteAtlas{
		Meta: Meta{
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
	"textures": [
		{
			"image": "{{.Config.ImageFile}}",
			"size": {
				"w": {{.Config.ImageWidth}},
				"h": {{.Config.ImageHeight}}
			},
			"sprites": [
				{{range .Rects}}
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
	],
	"meta": {
		"app": "{{.Mata}}"
	}
}`
