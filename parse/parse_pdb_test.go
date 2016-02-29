package parse

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParsePdb(t *testing.T) {
	file, _ := os.Open("../sample/hello.pdb")
	info, err := ParsePdb(file)
	if err != nil {
		t.Fatal(err)
	}

	// Debug Id
	assert.Equal(t, "B3963800A43840D28914331E6B93FE02", info.Guid.String())
	assert.Equal(t, 1, info.TimeDateStamp)
	assert.Equal(t, "B3963800A43840D28914331E6B93FE021", info.String())
}
