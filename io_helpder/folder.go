package io_helpder

import (
	"os"
	"path/filepath"
	"sort"
)

// ReadDirsRawList 阅读Dir 直接返回 类似于 os.ReadDir 但是不会排序
func ReadDirsRawList(dirPath string) ([]os.DirEntry, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.ReadDir(-1)
}

// ReadDirListRankByName 阅读dir 返回os.ReadDir 按照name排序
func ReadDirListRankByName(dirPath string, isAsc bool) ([]os.DirEntry, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dirs, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	sort.Slice(dirs, func(i, j int) bool {
		if isAsc {
			return dirs[i].Name() < dirs[j].Name()
		} else {
			return dirs[i].Name() > dirs[j].Name()
		}
	})
	return dirs, err
}

// ReadDirNames 阅读dir 返回names(短路径)
func ReadDirNames(dirPath string) ([]string, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)

	if err != nil {
		return nil, err
	}
	return names, nil

}

// ReadDirFullNames 阅读dir 返回names(长路径)
func ReadDirFullNames(dirPath string) ([]string, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)

	if err != nil {
		return nil, err
	}
	fullNames := make([]string, len(names))
	for i, name := range names {
		fullNames[i] = filepath.Join(absPath, name)
	}
	return fullNames, err
}
