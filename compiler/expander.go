package compiler

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Expander struct {
	rangeRegexp     *regexp.Regexp
	variationRegexp *regexp.Regexp
	hostnameRegexp  *regexp.Regexp
}

func NewExpander() *Expander {
	return &Expander{
		rangeRegexp:     regexp.MustCompile("\\[(\\d+)\\.\\.(\\d+)\\]"),
		variationRegexp: regexp.MustCompile("\\[([a-zA-Z0-9-|]+(?:\\.[a-zA-Z0-9-|]+)*)+\\]"),
		hostnameRegexp:  regexp.MustCompile("^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])(\\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]))*$"),
	}
}

type ExpandingRange struct {
	beginIdx int
	endIdx   int
	values   []string
}

type ExpandedHostname struct {
	Hostname     string
	Replacements []string
}

type ByIndex []ExpandingRange

func (s ByIndex) Len() int {
	return len(s)
}

func (s ByIndex) Less(i, j int) bool {
	return s[i].beginIdx < s[j].beginIdx
}

func (s ByIndex) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (e *Expander) expand(host string) ([]ExpandedHostname, error) {
	ranges := []ExpandingRange{}
	n := 1
	for _, r := range e.rangeRegexp.FindAllStringSubmatchIndex(host, -1) {
		expRange, err := e.expandingRange(host, r)
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, expRange)
		n *= len(expRange.values)
	}
	for _, v := range e.variationRegexp.FindAllStringSubmatchIndex(host, -1) {
		split := strings.Split(host[v[2]:v[3]], "|")
		ranges = append(ranges, ExpandingRange{
			beginIdx: v[0],
			endIdx:   v[1],
			values:   split,
		})
		n *= len(split)
	}
	if len(ranges) == 0 {
		return []ExpandedHostname{{Hostname: host}}, nil
	}
	hostnames, err := e.expandedHostnames(n, host, ranges)
	if err != nil {
		return nil, err
	}
	return hostnames, nil
}

func (e *Expander) expandingRange(host string, rangeGroup []int) (ExpandingRange, error) {
	begin, err := strconv.Atoi(host[rangeGroup[2]:rangeGroup[3]])
	if err != nil {
		return ExpandingRange{}, err
	}
	end, err := strconv.Atoi(host[rangeGroup[4]:rangeGroup[5]])
	if err != nil {
		return ExpandingRange{}, err
	}
	if begin >= end {
		return ExpandingRange{}, errors.New(fmt.Sprintf("Invalid range: %v is not smaller than %v", begin, end))
	}
	values := []string{}
	for i := begin; i <= end; i++ {
		values = append(values, strconv.Itoa(i))
	}
	return ExpandingRange{
		beginIdx: rangeGroup[0],
		endIdx:   rangeGroup[1],
		values:   values,
	}, nil
}

func (e *Expander) expandedHostnames(size int, host string, ranges []ExpandingRange) ([]ExpandedHostname, error) {
	hostnames := []ExpandedHostname{}
	sort.Sort(ByIndex(ranges))
	for i := 0; i < size; i++ {
		j := 1
		hostnameReplacements := []string{}
		produced := host[0:ranges[0].beginIdx]
		for p, r := range ranges {
			idx := (i / j) % len(r.values)
			value := r.values[idx]
			produced += value
			j *= len(r.values)
			nextIdx := p + 1
			if nextIdx < len(ranges) {
				produced += host[r.endIdx:ranges[nextIdx].beginIdx]
			} else {
				produced += host[r.endIdx:]
			}
			hostnameReplacements = append(hostnameReplacements, value)
		}
		if !e.hostnameRegexp.MatchString(produced) {
			return nil, errors.New(fmt.Sprintf("Produced string '%v' is not a valid Hostname", produced))
		}
		hostnames = append(hostnames, ExpandedHostname{
			Hostname:     produced,
			Replacements: hostnameReplacements,
		})
	}
	return hostnames, nil
}
