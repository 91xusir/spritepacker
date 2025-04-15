package pack

type AtlasInfo struct {
	Meta    Meta    `json:"meta"`
	Atlases []Atlas `json:"atlases"`
}

type Meta struct {
	Repo      string `json:"repo"`
	Format    string `json:"format"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type Atlas struct {
	Name    string   `json:"name"`
	Size    Size     `json:"size"`
	Sprites []Sprite `json:"sprites"`
}

type Sprite struct {
	FileName    string `json:"filename"`
	Frame       Rect   `json:"frame"`
	SrcRect     Size   `json:"srcRect"`
	TrimmedRect Rect   `json:"trimmedRect,omitzero"`
	Rotated     bool   `json:"rotated"`
	Trimmed     bool   `json:"trimmed"`
}

func (s Sprite) Clone() Sprite {
	return Sprite{
		FileName:    s.FileName,
		Frame:       s.Frame.Clone(),
		SrcRect:     s.SrcRect.Clone(),
		TrimmedRect: s.TrimmedRect.Clone(),
		Rotated:     s.Rotated,
		Trimmed:     s.Trimmed,
	}
}
