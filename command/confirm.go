package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Confirm struct {
	reader io.Reader
}

func NewConfirm(reader io.Reader) *Confirm {
	return &Confirm{
		reader: reader,
	}
}

func (c *Confirm) RequireConfirmationIfFileExists(path string) (bool, error) {
	exists, err := c.fileExists(path)
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}
	return c.confirmation(path)
}
func (c *Confirm) confirmation(path string) (bool, error) {
	r := bufio.NewReader(c.reader)
	fmt.Printf("File %s already exists. Overwrite? (Y/n)\n", path)
	response, err := r.ReadString('\n')
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(response) == "Y", nil
}

func (c *Confirm) fileExists(path string) (bool, error) {
	if strings.TrimSpace(path) == "" {
		return false, errors.New("Provided path is empty")
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if info.IsDir() {
		return true, errors.New(fmt.Sprintf("Path %s is a directory", path))
	}
	return true, nil
}
