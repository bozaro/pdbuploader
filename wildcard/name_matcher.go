package wildcard

import (
	"fmt"
	"regexp"
	"strings"
)

type NameMatcher interface {
	Matched(name string, dir bool) bool
	Recursive() bool
}

type maskType int

const (
	maskEquals maskType = iota
	maskSimple
	maskComplex
)

func NewNameMatcher(mask string) (NameMatcher, error) {
	if mask == "**/" {
		return RecursiveMatcher{}, nil
	}
	dirOnly := strings.HasSuffix(mask, "/")
	if dirOnly {
		mask = mask[:len(mask)-1]
	}
	nameMask := tryRemoveBackslashes(mask)
	reg, kind, err := maskToRegexp(nameMask)
	if err != nil {
		return nil, err
	}

	switch kind {
	case maskEquals:
		return EqualsMatcher{nameMask, dirOnly}, nil
	case maskSimple:
		asterisk := strings.Index(nameMask, "*")
		return SimpleMatcher{nameMask[:asterisk], nameMask[asterisk+1:], dirOnly}, nil
	default:
		return ComplexMatcher{reg, dirOnly}, nil
	}
}

func maskToRegexp(mask string) (*regexp.Regexp, maskType, error) {
	expr := ""
	last := 0
	kind := maskEquals
	for index, runeValue := range mask {
		switch runeValue {
		case '|', '-':
			expr += regexp.QuoteMeta(mask[last:index])
			expr += string(byte(runeValue))
			last = index + 1
		case '[', ']', '(', ')':
			expr += regexp.QuoteMeta(mask[last:index])
			expr += string(byte(runeValue))
			last = index + 1
			kind = maskComplex
		case '?':
			expr += regexp.QuoteMeta(mask[last:index])
			expr += "."
			last = index + 1
			kind = maskComplex
		case '*':
			expr += regexp.QuoteMeta(mask[last:index])
			expr += ".*"
			last = index + 1
			if kind == maskEquals {
				kind = maskSimple
			} else {
				kind = maskComplex
			}
		}
	}
	expr += regexp.QuoteMeta(mask[last:])
	reg, err := regexp.Compile(expr)
	return reg, kind, err
}

func tryRemoveBackslashes(pattern string) string {
	result := ""
	start := 0
	for true {
		next := strings.Index(pattern[start:], "\\")
		if next == -1 {
			if start < len(pattern) {
				result += pattern[start:]
			}
			break
		}
		next += start
		if next == len(pattern)-1 {
			// Return original string.
			return pattern
		}
		switch pattern[next+1] {
		case ' ', '#', '!':
			result += pattern[start:next]
			start = next + 1
			break
		default:
			return pattern
		}
	}
	return result
}

// Recursive directory matcher like "**".
type RecursiveMatcher struct {
}

func (this RecursiveMatcher) Matched(name string, dir bool) bool {
	return true
}

func (this RecursiveMatcher) Recursive() bool {
	return true
}

func (this RecursiveMatcher) String() string {
	return "**"
}

// Simple matcher for mask with only one asterisk.
type SimpleMatcher struct {
	prefix  string
	suffix  string
	dirOnly bool
}

func (this SimpleMatcher) Matched(name string, dir bool) bool {
	return (!this.dirOnly || dir) &&
		(len(name) >= len(this.prefix)+len(this.suffix)) &&
		strings.HasPrefix(name, this.prefix) &&
		strings.HasSuffix(name, this.suffix)
}

func (this SimpleMatcher) Recursive() bool {
	return false
}

func (this SimpleMatcher) String() string {
	return fmt.Sprintf("equals(%s, %s, %s)", this.prefix, this.suffix, this.dirOnly)
}

// Simple matcher for equals compare.
type EqualsMatcher struct {
	name    string
	dirOnly bool
}

func (this EqualsMatcher) Matched(name string, dir bool) bool {
	return (!this.dirOnly || dir) && (this.name == name)
}

func (this EqualsMatcher) Recursive() bool {
	return false
}

func (this EqualsMatcher) String() string {
	return fmt.Sprintf("equals(%s, %s)", this.name, this.dirOnly)
}

// Simple matcher for regexp compare.
type ComplexMatcher struct {
	matcher *regexp.Regexp
	dirOnly bool
}

func (this ComplexMatcher) Matched(name string, dir bool) bool {
	return (!this.dirOnly || dir) && this.matcher.MatchString(name)
}

func (this ComplexMatcher) Recursive() bool {
	return false
}

func (this ComplexMatcher) String() string {
	return fmt.Sprintf("complex(%s, %s)", this.matcher, this.dirOnly)
}
