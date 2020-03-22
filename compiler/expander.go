package compiler

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type expander struct {
	rangeRegexp     *regexp.Regexp
	variationRegexp *regexp.Regexp
	hostnameRegexp  *regexp.Regexp
}

func newExpander() *expander {
	return &expander{
		rangeRegexp:     regexp.MustCompile(`\[(\d+)\.\.(\d+)\]`),
		variationRegexp: regexp.MustCompile(`\[([a-zA-Z0-9-|]+(?:\.[a-zA-Z0-9-|]+)*)+\]`),
		hostnameRegexp:  regexp.MustCompile(`^([a-zA-Z0-9_]|[a-zA-Z0-9_][a-zA-Z0-9-_]{0,61}[a-zA-Z0-9_])(\.([a-zA-Z0-9_]|[a-zA-Z0-9_][a-zA-Z0-9-_]{0,61}[a-zA-Z0-9_]))*$`),
	}
}

type expandingRange struct {
	beginIdx int
	endIdx   int
	values   []string
}

type expandedHostname struct {
	Hostname     string
	Replacements []string
}

type byIndex []expandingRange

func (s byIndex) Len() int {
	return len(s)
}

func (s byIndex) Less(i, j int) bool {
	return s[i].beginIdx < s[j].beginIdx
}

func (s byIndex) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (e *expander) expand(host string) ([]expandedHostname, error) {
	var ranges []expandingRange
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
		ranges = append(ranges, expandingRange{
			beginIdx: v[0],
			endIdx:   v[1],
			values:   split,
		})
		n *= len(split)
	}
	if len(ranges) == 0 {
		if !e.hostnameRegexp.MatchString(host) {
			return nil, fmt.Errorf("produced string `%v` is not a valid Hostname", host)
		}
		return []expandedHostname{{Hostname: host}}, nil
	}
	hostnames, err := e.expandedHostnames(n, host, ranges)
	if err != nil {
		return nil, err
	}
	return hostnames, nil
}

func (e *expander) expandingRange(host string, rangeGroup []int) (expandingRange, error) {
	begin, err := strconv.Atoi(host[rangeGroup[2]:rangeGroup[3]])
	if err != nil {
		return expandingRange{}, err
	}
	end, err := strconv.Atoi(host[rangeGroup[4]:rangeGroup[5]])
	if err != nil {
		return expandingRange{}, err
	}
	if begin >= end {
		return expandingRange{}, fmt.Errorf("invalid range: %v is not smaller than %v", begin, end)
	}
	values := []string{}
	for i := begin; i <= end; i++ {
		values = append(values, strconv.Itoa(i))
	}
	return expandingRange{
		beginIdx: rangeGroup[0],
		endIdx:   rangeGroup[1],
		values:   values,
	}, nil
}

func (e *expander) expandedHostnames(size int, host string, ranges []expandingRange) ([]expandedHostname, error) {
	var hostnames []expandedHostname
	sort.Sort(byIndex(ranges))
	for i := 0; i < size; i++ {
		j := 1
		var hostnameReplacements []string
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
			return nil, fmt.Errorf("produced string `%v` is not a valid Hostname", produced)
		}
		hostnames = append(hostnames, expandedHostname{
			Hostname:     produced,
			Replacements: hostnameReplacements,
		})
	}
	return hostnames, nil
}
