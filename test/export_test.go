package export

import (
	"encoding/json"
	"github.com/91xusir/spritepacker/pack"
	"os"
	"testing"
)

func TestExport(t *testing.T) {
	// 测试导出功能
	// 打开中间结构 JSON 文件
	data, err := os.ReadFile("Atlas.json") // 假设你存的是这个文件
	if err != nil {
		t.Errorf("os.ReadFile failed: %v", err)
		return
	}

	// 反序列化为中间结构
	var atlas pack.AtlasInfo
	if err := json.Unmarshal(data, &atlas); err != nil {
		t.Errorf("json.Unmarshal failed: %v", err)
		return
	}

	var exportManager = NewExportManager()

	var godotExport Exporter = &GodotExporter{}

	exportManager.Register(godotExport.Ext(), godotExport)

	export, err := exportManager.Export("a.tpsheet", &atlas)
	if err != nil {
		t.Error(err)
	}
	err = os.WriteFile("a.tpsheet", export, 0644)
	if err != nil {
		t.Error(err)
	}
}
