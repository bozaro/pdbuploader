package wildcard

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type matchedResult int

const (
	matchedAlways matchedResult = iota
	matchedPossible
	matchedNever
)

type pathMatcherData struct {
	pattern       string
	path          string
	expectedMatch matchedResult
}

func (this pathMatcherData) String() string {
	return this.pattern + ", " + this.path
}

func TestPathMatcher(t *testing.T) {
	data := []pathMatcherData{
		// Simple pattern
		pathMatcherData{"/", "foo/bar", matchedAlways},
		pathMatcherData{"*", "foo/bar", matchedAlways},
		pathMatcherData{"*/", "foo/bar", matchedAlways},
		pathMatcherData{"/", "foo/bar/", matchedAlways},
		pathMatcherData{"*", "foo/bar/", matchedAlways},
		pathMatcherData{"*/", "foo/bar/", matchedAlways},
		pathMatcherData{"**/", "foo/bar/", matchedAlways},
		pathMatcherData{"foo/**/", "foo/bar/", matchedAlways},
		pathMatcherData{"foo/**/", "foo/bar/xxx", matchedAlways},
		pathMatcherData{"foo/**/", "foo/bar/xxx/", matchedAlways},
		pathMatcherData{"f*o", "foo/bar", matchedAlways},
		pathMatcherData{"/f*o", "foo/bar", matchedAlways},
		pathMatcherData{"f*o/", "foo/bar", matchedAlways},
		pathMatcherData{"foo/", "foo/bar", matchedAlways},
		pathMatcherData{"/foo/", "foo/bar", matchedAlways},
		pathMatcherData{"/foo", "foo/", matchedAlways},
		pathMatcherData{"foo", "foo/", matchedAlways},
		pathMatcherData{"foo/", "foo/", matchedAlways},
		pathMatcherData{"foo/", "foo", matchedPossible},
		pathMatcherData{"bar", "foo/bar", matchedAlways},
		pathMatcherData{"b*r", "foo/bar", matchedAlways},
		pathMatcherData{"/bar", "foo/bar", matchedNever},
		pathMatcherData{"bar/", "foo/bar", matchedPossible},
		pathMatcherData{"b*r/", "foo/bar", matchedPossible},
		pathMatcherData{"bar/", "foo/bar/", matchedAlways},
		pathMatcherData{"b*r/", "foo/bar/", matchedAlways},
		pathMatcherData{"b[a-z]r", "foo/bar", matchedAlways},
		pathMatcherData{"b[a-z]r", "foo/b0r", matchedPossible},
		pathMatcherData{"/t*e*t", "test", matchedAlways},
		// More complex pattern
		pathMatcherData{"foo/*/bar/", "foo/bar/", matchedPossible},
		pathMatcherData{"foo/*/bar/", "bar/", matchedNever},
		pathMatcherData{"foo/*/bar/", "foo/a/bar/", matchedAlways},
		pathMatcherData{"foo/*/bar/", "foo/a/b/bar/", matchedNever},
		pathMatcherData{"foo/*/*/bar/", "foo/a/b/bar/", matchedAlways},

		pathMatcherData{"foo/**/bar/a/", "foo/bar/b/bar/a/", matchedAlways},
		pathMatcherData{"foo/**/bar/a/", "foo/bar/bar/bar/a/", matchedAlways},
		pathMatcherData{"foo/**/bar/a/", "foo/bar/bar/b/a/", matchedPossible},
		pathMatcherData{"foo/**/bar/", "foo/bar/", matchedAlways},
		pathMatcherData{"foo/**/bar/", "bar/", matchedNever},
		pathMatcherData{"foo/**/bar/", "foo/a/bar/", matchedAlways},
		pathMatcherData{"foo/**/bar/", "foo/a/b/bar/", matchedAlways},
		pathMatcherData{"foo/*/**/*/bar/", "foo/a/bar/", matchedPossible},
		pathMatcherData{"foo/*/**/*/bar/", "foo/a/b/bar/", matchedAlways},
		pathMatcherData{"foo/*/**/*/bar/", "foo/a/b/c/bar/", matchedAlways},
		pathMatcherData{"foo/**/xxx/**/bar/", "foo/xxx/bar/", matchedAlways},
		pathMatcherData{"foo/**/xxx/**/bar/", "foo/xxx/b/c/bar/", matchedAlways},
		pathMatcherData{"foo/**/xxx/**/bar/", "foo/a/xxx/c/bar/", matchedAlways},
		pathMatcherData{"foo/**/xxx/**/bar/", "foo/a/c/xxx/bar/", matchedAlways},
		pathMatcherData{"foo/**/xxx/**/bar/", "foo/bar/xxx/", matchedPossible},
		pathMatcherData{"foo/**/xxx/**/bar/", "foo/bar/xxx/bar/", matchedAlways},
		pathMatcherData{"foo/**/xxx/**/bar/", "foo/bar/xxx/xxx/bar/", matchedAlways},
	}
	for _, value := range data {
		pathMatcherCheck(t, &value)
	}
}

func pathMatcherCheck(t *testing.T, v *pathMatcherData) {
	matcher, err := NewPathMatcher(v.pattern)
	assert.Nil(t, err, v.String())
	for _, name := range SplitPattern(v.path) {
		if matcher == nil {
			break
		}
		isDir := strings.HasSuffix(name, "/")
		if isDir {
			name = name[:len(name)-1]
		}
		matcher = matcher.CreateChild(name, isDir)
	}
	var actual matchedResult
	if matcher == nil {
		actual = matchedNever
	} else if matcher.Matched() {
		actual = matchedAlways
	} else {
		actual = matchedPossible
	}
	assert.Equal(t, v.expectedMatch, actual, v.String())
}

type tryRemoveBackslashesData struct {
	pattern  string
	expected string
}

func (this tryRemoveBackslashesData) String() string {
	return this.pattern
}

func TestTryRemoveBackslashes(t *testing.T) {
	data := []tryRemoveBackslashesData{
		tryRemoveBackslashesData{"test", "test"},
		tryRemoveBackslashesData{"test\\n", "test\\n"},
		tryRemoveBackslashesData{"space\\ ", "space "},
		tryRemoveBackslashesData{"foo\\!bar\\ ", "foo!bar "},
		tryRemoveBackslashesData{"\\#some", "#some"},
		tryRemoveBackslashesData{"foo\\[bar", "foo\\[bar"},
	}
	for _, value := range data {
		tryRemoveBackslashesCheck(t, &value)
	}
}

func tryRemoveBackslashesCheck(t *testing.T, v *tryRemoveBackslashesData) {
	actual := tryRemoveBackslashes(v.pattern)
	assert.Equal(t, v.expected, actual, v.String())
}
