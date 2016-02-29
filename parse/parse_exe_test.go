package parse

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseExe(t *testing.T) {
	file, _ := os.Open("../sample/hello.exe")
	info, err := ParseExe(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "C:\\Work\\github\\octobuild\\test_cl\\hello.pdb", info.PDBFileName)
	// Code Id
	assert.Equal(t, int32(0x56D047B1), info.CodeId.TimeDateStamp)
	assert.Equal(t, int32(0x2e000), info.CodeId.SizeOfImage)
	assert.Equal(t, "56D047B12e000", info.CodeId.String())
	// Debug Id
	assert.Equal(t, "B3963800A43840D28914331E6B93FE02", info.DebugId.Guid.String())
	assert.Equal(t, 1, info.DebugId.TimeDateStamp)
	assert.Equal(t, "B3963800A43840D28914331E6B93FE021", info.DebugId.String())
}
