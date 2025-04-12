package pack

type SpriteAtlas struct {
	Meta    Meta    `json:"meta"`
	Atlases []Atlas `json:"atlases"` // 改名为复数形式，更符合内容
}

type Meta struct {
	Repo      string `json:"repo"`
	Format    string `json:"format"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type Atlas struct {
	Name    string   `json:"name"`
	Size    Size     `json:"size"` // 修正为小写，与其他字段保持一致
	Sprites []Sprite `json:"sprites"`
}

type Sprite struct {
	Filepath    string    `json:"filepath"`
	Frame       Rectangle `json:"frame"`
	SrcRect     Size      `json:"srcRect"`
	TrimmedRect Rectangle `json:"trimmedRect,omitempty"`
	Rotated     bool      `json:"rotated"`
	Trimmed     bool      `json:"trimmed"`
}

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Rectangle struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}
