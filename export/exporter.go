package export

import (
	"encoding/json"
	"github.com/91xusir/spritepacker/pack"
	"path/filepath"
	"strings"
)

type Exporter interface {
	Export(atlas *pack.SpriteAtlas) ([]byte, error)
	Import(data []byte) (*pack.SpriteAtlas, error)
}

type ExporterManager struct {
	exporters map[string]Exporter
}

func NewExportManager() *ExporterManager {
	return &ExporterManager{
		exporters: make(map[string]Exporter),
	}
}

func (m *ExporterManager) Register(ext string, exporter Exporter) {
	ext = strings.ToLower(ext)
	m.exporters[ext] = exporter
}

func (m *ExporterManager) Export(fileName string, atlas *pack.SpriteAtlas) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	if exporter, ok := m.exporters[ext]; ok {
		return exporter.Export(atlas)
	}
	return json.MarshalIndent(atlas, "", "  ") // 默认JSON
}

func (m *ExporterManager) Import(fileName string, data []byte) (*pack.SpriteAtlas, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	if exporter, ok := m.exporters[ext]; ok {
		return exporter.Import(data)
	}
	var atlas pack.SpriteAtlas
	err := json.Unmarshal(data, &atlas)
	return &atlas, err
}
