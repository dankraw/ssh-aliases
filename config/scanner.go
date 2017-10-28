package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Scanner struct{}

func NewScanner() *Scanner {
	return &Scanner{}
}

const hclExtension = ".hcl"

func (s *Scanner) ScanDirectory(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	hcls := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), hclExtension) {
			hcls = append(hcls, filepath.Join(path, file.Name()))
		}
	}
	return hcls, nil
}

func (s *Scanner) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
