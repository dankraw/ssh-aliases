package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Scanner is used to select files that contain ssh-aliases configs
type Scanner struct{}

// NewScanner creates new instance of Scanner
func NewScanner() *Scanner {
	return &Scanner{}
}

const hclExtension = ".hcl"

// ScanDirectory returns an array of file names that contain ssh-aliases configs
func (s *Scanner) ScanDirectory(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error while scanning `%s`: %s", path, err.Error())
	}
	var hcls []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), hclExtension) {
			hcls = append(hcls, filepath.Join(path, file.Name()))
		}
	}
	return hcls, nil
}
