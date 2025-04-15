package pack

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// This file provides utility functions for image and file operations.
// Image operations are based on the https://github.com/disintegration/imaging package.
// Some implementations were translated with the help of https://chat.deepseek.com.
// And a few are my own creations, haha.

//---------------Math----------------

// MaxInt returns the larger of a and b.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt returns the smaller of a and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// NextPowerOfTwo returns the next power of two of n.
func NextPowerOfTwo(n int) int {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

//---------------file--------------------

func SafeCreate(outputPath string) (*os.File, error) {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(outputPath)
}

// ListFilePaths
//
// Parameters:
//   - dirPath: the directory path
//
// Returns:
//   - []string: the file paths
//   - error: the error
//
// Example:
//
//	filePaths, err := ListFilePaths("./sprites")
func ListFilePaths(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, errors.New("failed to read directory: " + err.Error())
	}
	var paths = make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			path := filepath.Join(dirPath, entry.Name())
			paths = append(paths, path)
		}
	}
	// sort file names in natural order
	NaturalSort(paths)
	return paths, nil
}

// GetLastFolderName
//
// Parameters:
//   - path: the path
//
// Returns:
//   - string: the last folder name
//
// Example:
//
//	folderName := GetLastFolderName("./sprites/atlas/a.png")
//	// folderName = "atlas"
func GetLastFolderName(path string) string {
	path = filepath.ToSlash(path)
	path = strings.TrimRight(path, "/")
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return "atlas"
}

var chunkedRegex = regexp.MustCompile(`(\d+|\D+)`)

func chunked(path string) []string {
	return chunkedRegex.FindAllString(path, -1)
}

func NaturalLess(a, b string) bool {
	chunksA := chunked(a)
	chunksB := chunked(b)
	for i := 0; i < len(chunksA) && i < len(chunksB); i++ {
		aChunk := chunksA[i]
		bChunk := chunksB[i]
		if aChunk != bChunk {
			aNum, aErr := strconv.Atoi(aChunk)
			bNum, bErr := strconv.Atoi(bChunk)
			if aErr == nil && bErr == nil {
				return aNum < bNum
			}
			return aChunk < bChunk
		}
	}
	return len(chunksA) < len(chunksB)
}

// NaturalSort sorts the items in natural order.
func NaturalSort(items []string) {
	sort.Slice(items, func(i, j int) bool {
		return NaturalLess(items[i], items[j])
	})
}

type SameDetectInfo struct {
	DupeToBaseName  map[string]string
	BaseToDupesName map[string][]string
}

// FindDuplicateFiles finds duplicate files in the given file paths.
//
// Parameters:
//   - filePaths: the paths of the files to be checked
//
// Returns:
//   - []string: the paths of the unique files
//   - map[string]string: the map of duplicate file names and their corresponding base names
//   - map[string][]string: the map of base file names to their duplicate file names
//   - error
func FindDuplicateFiles(filePaths []string) ([]string, SameDetectInfo, error) {
	var uniqueFiles []string
	reverseDupeMap := make(map[string]string)
	// 新增映射：源文件对应重复文件的数组
	duplicatesMap := make(map[string][]string)
	filesBySize := make(map[int64][]string)
	for _, path := range filePaths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, SameDetectInfo{}, err
		}

		if !info.IsDir() {
			filesBySize[info.Size()] = append(filesBySize[info.Size()], path)
		}
	}
	for _, files := range filesBySize {
		if len(files) == 1 {
			uniqueFiles = append(uniqueFiles, files[0])
			continue
		}
		filesByHash := make(map[string][]string)
		for _, file := range files {
			hash, err := calculateMD5(file)
			if err != nil {
				continue
			}
			filesByHash[hash] = append(filesByHash[hash], file)
		}
		for _, paths := range filesByHash {
			if len(paths) == 1 {
				uniqueFiles = append(uniqueFiles, paths[0])
			} else {
				uniqueFiles = append(uniqueFiles, paths[0])
				keepFile := paths[0]
				baseKeepFile := filepath.Base(keepFile)
				duplicatesMap[baseKeepFile] = []string{}
				for _, dupFile := range paths[1:] {
					baseDupFile := filepath.Base(dupFile)
					reverseDupeMap[baseDupFile] = baseKeepFile
					duplicatesMap[baseKeepFile] = append(duplicatesMap[baseKeepFile], baseDupFile)
				}
			}
		}
	}
	NaturalSort(uniqueFiles)
	return uniqueFiles, SameDetectInfo{
		DupeToBaseName:  reverseDupeMap,
		BaseToDupesName: duplicatesMap,
	}, nil
}

func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

//---------------other--------------------

// Parallel  processes the data in separate goroutines.
func Parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}
	process := runtime.GOMAXPROCS(0)
	if process > count {
		process = count
	}

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)
	var wg sync.WaitGroup
	for i := 0; i < process; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}
