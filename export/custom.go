package export

import (
	"bytes"
	"errors"
	"github.com/91xusir/spritepacker/model"
	"text/template"
)

// TemplateExporter 模板导出器
// 提供了一个通用的导出器，可以通过模板字符串来导出数据
// 模板字符串中可以使用 {{.Atlas}} 来访问 AtlasInfo
type TemplateExporter struct {
	ext          string
	templateStr  string
	parseFunc    func([]byte) (*model.AtlasInfo, error)
	templateFunc template.FuncMap
}
type ParseFunc func([]byte) (*model.AtlasInfo, error)

func NewTemplateExporter(templateStr string, parseFunc ParseFunc) *TemplateExporter {
	return &TemplateExporter{
		ext:         ".tmpl",
		templateStr: templateStr,
		parseFunc:   parseFunc,
		templateFunc: template.FuncMap{
			"isLast": func(index, length int) bool {
				return index == length-1
			},
		},
	}
}

func (e *TemplateExporter) Export(atlas *model.AtlasInfo) ([]byte, error) {
	if e.templateStr == "" {
		return nil, errors.New("template string is empty")
	}

	tmpl, err := template.New("export").Funcs(e.templateFunc).Parse(e.templateStr)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, atlas); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *TemplateExporter) Import(data []byte) (*model.AtlasInfo, error) {
	if e.parseFunc == nil {
		return nil, errors.New("parse function is not provided")
	}
	return e.parseFunc(data)
}

func (e *TemplateExporter) Ext() string {
	return e.ext
}

func (e *TemplateExporter) SetExt(ext string) {
	e.ext = ext
}
func (e *TemplateExporter) AddTemplateFunc(name string, fn interface{}) {
	if e.templateFunc == nil {
		e.templateFunc = template.FuncMap{}
	}
	e.templateFunc[name] = fn
}
