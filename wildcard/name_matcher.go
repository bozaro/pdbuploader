package wildcard

import (
	"regexp"
	"strings"
)

type NameMatcher interface {
	Matched(name string, dir bool) bool
	Recursive() bool
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
