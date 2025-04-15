package export

import (
	"encoding/json"
	"github.com/91xusir/spritepacker/pack"
)

type JsonExporter struct {
	ext string
}

func (j *JsonExporter) Ext() string {
	return j.ext
}
func (j *JsonExporter) SetExt(ext string) {
	j.ext = ext
}

func (j *JsonExporter) Export(atlas *pack.AtlasInfo) ([]byte, error) {
	return json.Marshal(atlas)
}

func (j *JsonExporter) Import(data []byte) (*pack.AtlasInfo, error) {
	var atlas pack.AtlasInfo
	err := json.Unmarshal(data, &atlas)
	return &atlas, err
}
