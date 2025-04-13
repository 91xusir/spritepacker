package pack

import "testing"

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
	NaturalSort(files)
	for i, file := range files {
		t.Logf("file[%d] = %s", i, file)
	}
}
