package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type confirm struct {
	reader io.Reader
}

func newConfirm(reader io.Reader) *confirm {
	return &confirm{
		reader: reader,
	}
}

func (c *confirm) requireConfirmationIfFileExists(path string) (bool, error) {
	exists, err := c.fileExists(path)
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}
	return c.confirmation(path)
}
func (c *confirm) confirmation(path string) (bool, error) {
	r := bufio.NewReader(c.reader)
	fmt.Printf("File `%s` already exists. Overwrite? (Y/n)\n", path)
	response, err := r.ReadString('\n')
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(response) == "Y", nil
}

func (c *confirm) fileExists(path string) (bool, error) {
	if strings.TrimSpace(path) == "" {
		return false, errors.New("provided path is empty")
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if info.IsDir() {
		return true, fmt.Errorf("path `%s` is a directory", path)
	}
	return true, nil
}
