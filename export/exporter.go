package export

import (
	"encoding/json"
	"errors"
	"github.com/91xusir/spritepacker/pack"
	"os"
	"path/filepath"
	"strings"
)

type Exporter interface {
	Export(atlas *pack.AtlasInfo) ([]byte, error)
	Import(data []byte) (*pack.AtlasInfo, error)
	Ext() string
	SetExt(ext string)
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
	exporter.SetExt(ext)
	m.exporters[ext] = exporter
}

func (m *ExporterManager) Export(fileName string, atlas *pack.AtlasInfo) error {
	ext := strings.ToLower(filepath.Ext(fileName))
	var data []byte
	var err error
	if exporter, ok := m.exporters[ext]; ok {
		data, err = exporter.Export(atlas)
		if err != nil {
			return err
		}
	} else {
		return errors.New("unsupported file type")
	}
	if err = os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func (m *ExporterManager) Import(fileName string) (*pack.AtlasInfo, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	if exporter, ok := m.exporters[ext]; ok {
		return exporter.Import(data)
	}
	var atlas pack.AtlasInfo
	err = json.Unmarshal(data, &atlas)
	return &atlas, err
}

func (m *ExporterManager) Init() *ExporterManager {
	m.Register(".json", &JsonExporter{})
	m.Register(".tpsheet", &GodotExporter{})
	return m
}

func (m *ExporterManager) RegisterTemplate(ext string, templateStr string, parseFunc ParseFunc) {
	if parseFunc == nil {
		parseFunc = func(data []byte) (*pack.AtlasInfo, error) {
			return nil, errors.New("parse function is not provided")
		}
	}
	exporter := NewTemplateExporter(templateStr, parseFunc)
	m.Register(ext, exporter)
}
