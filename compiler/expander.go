package compiler

import (
	"regexp"
	"fmt"
	"strconv"
	"errors"
	"strings"
	"sort"
)

type Expander struct {
	rangeRegexp *regexp.Regexp
	variationRegexp *regexp.Regexp
}

func NewExpander() *Expander {
	return &Expander{
		rangeRegexp:     regexp.MustCompile("\\[(\\d+)\\.\\.(\\d+)\\]"),
		variationRegexp: regexp.MustCompile("\\[([a-zA-Z0-9-|]+(?:\\.[a-zA-Z0-9-|]+)*)+\\]"),
	}
}

type stringGroup struct {
	beginIdx int
	endIdx int
	values []string
}

type ByIndex []stringGroup

func (s ByIndex) Len() int {
	return len(s)
}

func (s ByIndex) Less(i, j int) bool {
	return s[i].beginIdx < s[j].beginIdx
}

func (s ByIndex) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (e *Expander) expand(host string) ([]string, error) {
	groups := []stringGroup{}
	n := 1
	for _, r := range e.rangeRegexp.FindAllStringSubmatchIndex(host, -1) {
		begin, err := strconv.Atoi(host[r[2]:r[3]])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(host[r[4]:r[5]])
		if err != nil {
			return nil, err
		}
		if begin >= end {
			return nil, errors.New(fmt.Sprintf("Invalid range: %v is not smaller than %v", begin, end))
		}
		rArray := []string{}
		for i:= begin; i <= end; i++ {
			rArray = append(rArray, strconv.Itoa(i))
		}
		groups = append(groups, stringGroup{
			beginIdx: r[0],
			endIdx: r[1],
			values: rArray,
		})
		n *= len(rArray)
	}
	for _, v := range e.variationRegexp.FindAllStringSubmatchIndex(host, -1) {
		split := strings.Split(host[v[2]:v[3]], "|")
		groups = append(groups, stringGroup{
			beginIdx: v[0],
			endIdx:   v[1],
			values:   split,
		})
		n *= len(split)
	}
	if len(groups) == 0 {
		return []string{host}, nil
	}
	sort.Sort(ByIndex(groups))

	hostnames := []string{}
	for i := 0; i < n; i++ {
		j := 1
		produced := host[0:groups[0].beginIdx]
		for p, r := range groups {
			idx := (i / j) % len(r.values)
			produced += r.values[idx]
			j *= len(r.values)
			if p < len(groups) - 1 {
				produced += host[r.endIdx:groups[p+1].beginIdx]
			} else {
				produced += host[r.endIdx:]
			}
		}
		hostnames = append(hostnames, produced)
	}
	return hostnames, nil
}
