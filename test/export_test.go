package spritepacker

import (
	"encoding/json"
	"github.com/91xusir/spritepacker/export"
	"github.com/91xusir/spritepacker/pack"
	"os"
	"testing"
	"time"
)

func TestExport(t *testing.T) {
	data, err := os.ReadFile("./output/atlas.json")
	if err != nil {
		t.Errorf("os.ReadFile failed: %v", err)
		return
	}
	var atlas pack.AtlasInfo
	if err = json.Unmarshal(data, &atlas); err != nil {
		t.Errorf("json.Unmarshal failed: %v", err)
		return
	}

	var exportManager = export.NewExportManager().Init()

	err = exportManager.Export("atlas.tpsheet", &atlas)
	if err != nil {
		t.Errorf("exportManager.Export failed: %v", err)
	}

}
func TestTmpl(t *testing.T) {
	manager := export.NewExportManager().Init()
	jsTemplate := `
const ATLAS_DATA = {
  meta: {
    repo: "{{.Meta.Repo}}",
    format: "{{.Meta.Format}}",
    version: "{{.Meta.Version}}",
    timestamp: "{{.Meta.Timestamp}}"
  },
  atlases: [
    {{- range $i, $atlas := .Atlases }}
    {
      name: "{{$atlas.Name}}",
      size: { w: {{$atlas.Size.W}}, h: {{$atlas.Size.H}} },
      sprites: [
        {{- range $j, $sprite := $atlas.Sprites }}
        {
          filename: "{{$sprite.FileName}}",
          frame: { x: {{$sprite.Frame.X}}, y: {{$sprite.Frame.Y}}, w: {{$sprite.Frame.W}}, h: {{$sprite.Frame.H}} },
          rotated: {{$sprite.Rotated}},
          trimmed: {{$sprite.Trimmed}}
        }{{if not (isLast $j (len $atlas.Sprites))}},{{end}}
        {{- end }}
      ]
    }{{if not (isLast $i (len $.Atlases))}},{{end}}
    {{- end }}
  ]
};
export default ATLAS_DATA;
`
	manager.RegisterTemplate(".js", jsTemplate, nil)
	atlas := &pack.AtlasInfo{
		Meta: pack.Meta{
			Repo:      "https://github.com/user/spritepacker",
			Format:    "RGBA8888",
			Version:   "1.0.0",
			Timestamp: time.Now().Format(time.RFC3339),
		},
		Atlases: []pack.Atlas{
			{
				Name: "main",
				Size: pack.Size{W: 1024, H: 1024},
				Sprites: []pack.Sprite{
					{
						FileName: "player.png",
						Frame:    pack.NewRectByPosAndSize(0, 0, 64, 64),
						SrcRect:  pack.Size{W: 64, H: 64},
						Rotated:  false,
						Trimmed:  false,
					},
					{
						FileName: "enemy.png",
						Frame:    pack.NewRectByPosAndSize(0, 0, 64, 64),
						SrcRect:  pack.Size{W: 32, H: 32},
						Rotated:  false,
						Trimmed:  false,
					},
				},
			},
		},
	}
	jsFilePath := "./atlas.js"
	if err := manager.Export(jsFilePath, atlas); err != nil {
		t.Fatalf("%v", err)
	}
}
