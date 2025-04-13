package pack

type SpriteAtlas struct {
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
	Filepath    string    `json:"filepath"`
	Frame       Rectangle `json:"frame"`
	SrcRect     Size      `json:"srcRect"`
	TrimmedRect Rectangle `json:"trimmedRect,omitzero"`
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
