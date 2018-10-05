package command

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func readHostsFile(hostsFile string) ([]string, error) {
	if hostsFile == "" {
		return []string{}, nil
	}
	bytes, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read input hosts file: %s: %s", hostsFile, err.Error())
	}
	return strings.Split(string(bytes), "\n"), nil
}
