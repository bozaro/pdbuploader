package wildcard

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type splitPatternData struct {
	pattern  string
	expected []string
}

func (this splitPatternData) String() string {
	return this.pattern
}

func TestSplitPattern(t *testing.T) {
	data := []splitPatternData{
		splitPatternData{"foo", []string{"foo"}},
		splitPatternData{"foo/", []string{"foo/"}},
		splitPatternData{"/bar", []string{"/", "bar"}},
		splitPatternData{"/foo/bar/**", []string{"/", "foo/", "bar/", "**"}},
	}
	for _, value := range data {
		splitPatternCheck(t, &value)
	}
}

func splitPatternCheck(t *testing.T, v *splitPatternData) {
	actual := SplitPattern(v.pattern)
	assert.Equal(t, v.expected, actual, v.String())
}
