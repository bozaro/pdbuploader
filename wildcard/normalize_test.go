package wildcard

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type normalizePatternData struct {
	pattern  string
	expected []string
}

func (this normalizePatternData) String() string {
	return this.pattern
}

func TestNormalizePattern(t *testing.T) {
	data := []normalizePatternData{
		// Simple mask
		normalizePatternData{"/", []string{}},
		normalizePatternData{"*/", []string{"*/"}},
		normalizePatternData{"*", []string{}},
		normalizePatternData{"**", []string{}},
		normalizePatternData{"**/", []string{}},
		normalizePatternData{"foo", []string{"**/", "foo"}},
		normalizePatternData{"foo/", []string{"**/", "foo/"}},
		normalizePatternData{"/foo", []string{"foo"}},

		// Convert path file mask
		normalizePatternData{"foo/**.bar", []string{"foo/", "**/", "*.bar"}},
		normalizePatternData{"foo/***.bar", []string{"foo/", "**/", "*.bar"}},

		// Collapse and reorder adjacent masks
		normalizePatternData{"foo/*/bar", []string{"foo/", "*/", "bar"}},
		normalizePatternData{"foo/**/bar", []string{"foo/", "**/", "bar"}},
		normalizePatternData{"foo/*/*/bar", []string{"foo/", "*/", "*/", "bar"}},
		normalizePatternData{"foo/**/*/bar", []string{"foo/", "*/", "**/", "bar"}},
		normalizePatternData{"foo/*/**/bar", []string{"foo/", "*/", "**/", "bar"}},
		normalizePatternData{"foo/*/**.bar", []string{"foo/", "*/", "**/", "*.bar"}},
		normalizePatternData{"foo/**/**/bar", []string{"foo/", "**/", "bar"}},
		normalizePatternData{"foo/**/**.bar", []string{"foo/", "**/", "*.bar"}},
		normalizePatternData{"foo/**/*/**/*/bar", []string{"foo/", "*/", "*/", "**/", "bar"}},
		normalizePatternData{"foo/**/*/**/*/**.bar", []string{"foo/", "*/", "*/", "**/", "*.bar"}},

		// Collapse trailing masks
		normalizePatternData{"foo/**", []string{"foo/"}},
		normalizePatternData{"foo/**/*", []string{"foo/"}},
		normalizePatternData{"foo/**/*/*", []string{"foo/", "*/"}},
		normalizePatternData{"foo/**/", []string{"foo/"}},
		normalizePatternData{"foo/**/*/", []string{"foo/", "*/"}},
		normalizePatternData{"foo/**/*/*/", []string{"foo/", "*/", "*/"}},
	}
	for _, value := range data {
		normalizePatternCheck(t, &value)
	}
}

func normalizePatternCheck(t *testing.T, v *normalizePatternData) {
	actual := NormalizePattern(SplitPattern(v.pattern))
	assert.Equal(t, append([]string{"/"}, v.expected...), actual, v.String())
}
