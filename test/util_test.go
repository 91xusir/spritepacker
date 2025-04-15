package spritepacker

import (
	"github.com/91xusir/spritepacker/utils"
	"testing"
)

func TestNaturalSort(t *testing.T) {
	files := []string{
		"file1.txt",
		"file10.txt",
		"file2.txt",
		"file20.txt",
		"file3.txt",
		"file30.txt",
		"/path/to/file1.txt",
		"/path/to/file10.txt",
		"/path/to/file2.txt",
	}
	utils.NaturalSort(files)
	for i, file := range files {
		t.Logf("file[%d] = %s", i, file)
	}
}

// util_test.go:25: before len 65
// util_test.go:27: after len 62
// util_test.go:29: key 30.png, value 29.png
// util_test.go:29: key 31.png, value 29.png
// util_test.go:29: key 32.png, value 29.png
// util_test.go:32: key 29.png, value [30.png 31.png 32.png]
func TestFindDuplicateImages(t *testing.T) {
	files, _ := utils.ListFilePaths("../test/input")
	t.Logf("before len %d ", len(files))
	paths, info, _ := utils.FindDuplicateFiles(files)
	t.Logf("after len %d ", len(paths))
	for k, v := range info.DupeToBaseName {
		t.Logf("key %s, value %v", k, v)
	}
	for k, v := range info.BaseToDupesName {
		t.Logf("key %s, value %v", k, v)
	}
}
