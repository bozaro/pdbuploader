package parse

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParsePdb32(t *testing.T) {
	file, _ := os.Open("../sample/hello32.pdb")
	info, err := ParsePdb(file)
	if err != nil {
		t.Fatal(err)
	}

	// Debug Id
	assert.Equal(t, "74219E6FC19941AAB1812BA7E7FF1EFE", info.Guid.String())
	assert.Equal(t, 2, info.TimeDateStamp)
	assert.Equal(t, "74219E6FC19941AAB1812BA7E7FF1EFE2", info.String())
}

func TestParsePdb64(t *testing.T) {
	file, _ := os.Open("../sample/hello64.pdb")
	info, err := ParsePdb(file)
	if err != nil {
		t.Fatal(err)
	}

	// Debug Id
	assert.Equal(t, "B39838386DAE42DF9F1E6B0B12963079", info.Guid.String())
	assert.Equal(t, 2, info.TimeDateStamp)
	assert.Equal(t, "B39838386DAE42DF9F1E6B0B129630792", info.String())
}
