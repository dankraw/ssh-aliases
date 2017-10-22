package compiler

import (
	"regexp"
	"fmt"
	"strconv"
	"errors"
)

type Expander struct {
	rangeRegexp *regexp.Regexp
}

func NewExpander() *Expander {
	rangeRegexp, _ := regexp.Compile("\\[(\\d+)\\.\\.(\\d+)\\]")
	return &Expander{
		rangeRegexp: rangeRegexp,
	}
}

func (e *Expander) expand(host string) ([]string, error) {
	hostnames := []string{host}
	tryExpanding := true
	for tryExpanding {
		nextIterationHostnames := []string{}
		for i := 0; i < len(hostnames); i++ {
			expanded, err := e.expandWithNextRangeFound(hostnames[i])
			if err != nil {
				return nil, err
			}
			nextIterationHostnames = append(nextIterationHostnames, expanded...)
		}
		if len(hostnames) >= len(nextIterationHostnames) {
			tryExpanding = false
		}
		hostnames = nextIterationHostnames
	}
	return hostnames, nil
}

func (e *Expander) expandWithNextRangeFound(host string) ([]string, error) {
	idx := e.rangeRegexp.FindAllStringSubmatchIndex(host, -1)
	for _, group := range idx {
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
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid range: %v is not smaller than %v", begin, end))
		}
	}
	return []string{host}, nil
}

