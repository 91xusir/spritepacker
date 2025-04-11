package spritepacker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"spritepacker/internal/packer"
	"spritepacker/internal/shape"
	"strconv"
	"strings"
	"testing"
)

func Test_Packing(t *testing.T) {
	// you can use func NewOptions to set options
	opt := packer.NewOptions(
		packer.WithMaxSize(1024, 1024),
		packer.WithPadding(0),
		packer.WithAlgorithm(packer.AlgoSkyline),
	)
	_ = opt
	// you also can use OptionBuilder to set options
	opt = packer.NewOptionBuilder().MaxSize(1024, 1024).Padding(0).Algorithm(packer.AlgoSkyline).Build()


	reqPackedRect, err := getTestData("testdata/rect.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	bin := shape.NewBin(400, 400, reqPackedRect)

	packer := packer.NewPacker(bin, opt)
	packResult := packer.Packing()
	fmt.Println(packResult)

}

func getTestData(path string) ([]*shape.Rect, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("cannot get current file path")
	}
	base := filepath.Dir(filename)
	absPath := filepath.Join(base, path)
	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var reqPackedRect []*shape.Rect
	id := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		w, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("an error in parsing rectangle width: %w", err)
		}
		h, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("an error in parsing rectangle height: %w", err)
		}
		rect := shape.NewRectById(id, w, h)
		id++
		reqPackedRect = append(reqPackedRect, rect)

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return reqPackedRect, nil
}
