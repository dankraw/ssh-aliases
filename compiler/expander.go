package compiler

import (
	"regexp"
	"fmt"
	"strconv"
	"errors"
	"strings"
)

type Expander struct {
	rangeRegexp *regexp.Regexp
	variationRegexp *regexp.Regexp
}

func NewExpander() *Expander {
	return &Expander{
		rangeRegexp:     regexp.MustCompile("\\[(\\d+)\\.\\.(\\d+)\\]"),
		variationRegexp: regexp.MustCompile("\\[([a-zA-Z0-9-.|]+)\\]"),
	}
}

func (e *Expander) expand(host string) ([]string, error) {
	hostnames, err := e.expandWithFunction([]string{host}, e.expandWithNextRangeFound)
	if err != nil {
		return nil, err
	}
	hostnames, err = e.expandWithFunction(hostnames, e.expandWithNextVariationFound)
	if err != nil {
		return nil, err
	}
	return hostnames, nil
}

func (e *Expander) expandWithFunction(hostnames []string, expandWith func (string) ([]string, error)) ([]string, error) {
	tryExpanding := true
	expanded := hostnames
	for tryExpanding {
		nextIterationHostnames := []string{}
		for i := 0; i < len(expanded); i++ {
			expanded, err := expandWith(expanded[i])
			if err != nil {
				return nil, err
			}
			nextIterationHostnames = append(nextIterationHostnames, expanded...)
		}
		if len(expanded) >= len(nextIterationHostnames) {
			tryExpanding = false
		}
		expanded = nextIterationHostnames
	}
	return expanded, nil
}

func (e *Expander) expandWithNextRangeFound(host string) ([]string, error) {
	group := e.rangeRegexp.FindStringSubmatchIndex(host)
	if len(group) < 1 {
		return []string{host}, nil
	}
	begin, err := strconv.Atoi(host[group[2]:group[3]])
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(host[group[4]:group[5]])
	if err != nil {
		return nil, err
	}
	if begin < end {
		hostnames := []string{}
		for i := begin; i <= end; i++ {
			expanded := fmt.Sprintf("%s%v%s", host[0:group[0]], i, host[group[1]:])
			hostnames = append(hostnames, expanded)
		}
		return hostnames, nil
	}
	return nil, errors.New(fmt.Sprintf("Invalid range: %v is not smaller than %v", begin, end))
}

func (e *Expander) expandWithNextVariationFound(host string) ([]string, error) {
	group := e.variationRegexp.FindStringSubmatchIndex(host)
	if len(group) < 1 {
		return []string{host}, nil
	}
	variations := strings.Split(host[group[2]:group[3]], "|")
	hostnames := []string{}
	for _, v := range variations {
		expanded := fmt.Sprintf("%s%v%s", host[0:group[0]], v, host[group[1]:])
		hostnames = append(hostnames, expanded)
	}
	return hostnames, nil
}

