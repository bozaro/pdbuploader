package wildcard

import (
	"regexp"
	"strings"
)

type NameMatcher interface {
	Matched(name string, dir bool) bool
	Recursive() bool
}

func NewNameMatcher(mask string) NameMatcher {
	if mask == "**/" {
		return RecursiveMatcher{}
	}
	dirOnly := strings.HasSuffix(mask, "/")
	if dirOnly {
		mask = mask[:len(mask)-1]
	}
	nameMask := tryRemoveBackslashes(mask)
	if strings.Contains(nameMask, "[") || strings.Contains(nameMask, "]") || strings.Contains(nameMask, "\\") || strings.Contains(nameMask, "?") {
		return ComplexMatcher{maskToRegexp(nameMask), dirOnly}
	}
	// Subversion compatible mask.
	asterisk := strings.Index(nameMask, "*")
	if asterisk < 0 {
		return EqualsMatcher{nameMask, dirOnly}
	} else if strings.Index(mask[asterisk+1:], "*") < 0 {
		return SimpleMatcher{nameMask[:asterisk], nameMask[asterisk+1:], dirOnly}
	}
	return ComplexMatcher{maskToRegexp(nameMask), dirOnly}
}

func maskToRegexp(mask string) *regexp.Regexp {
	// todo: Mask to regexp
	return nil
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
		case ' ':
		case '#':
		case '!':
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
	return dir
}

func (this RecursiveMatcher) Recursive() bool {
	return true
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
