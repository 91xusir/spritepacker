package export

import (
	"encoding/json"
	"github.com/91xusir/spritepacker/model"
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

func (j *JsonExporter) Export(atlas *model.AtlasInfo) ([]byte, error) {
	return json.MarshalIndent(atlas, "", "    ")
}

func (j *JsonExporter) Import(data []byte) (*model.AtlasInfo, error) {
	var atlas model.AtlasInfo
	err := json.Unmarshal(data, &atlas)
	return &atlas, err
}
