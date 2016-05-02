package parse

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseExe32(t *testing.T) {
	file, _ := os.Open("../sample/hello32.exe")
	info, err := ParseExe(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "C:\\Bozaro\\pdbuploader\\sample\\hello32.pdb", info.PDBFileName)
	// Code Id
	assert.Equal(t, int32(0x57218B93), info.CodeId.TimeDateStamp)
	assert.Equal(t, int32(0x65000), info.CodeId.SizeOfImage)
	assert.Equal(t, "57218B9365000", info.CodeId.String())
	// Debug Id
	assert.Equal(t, "74219E6FC19941AAB1812BA7E7FF1EFE", info.DebugId.Guid.String())
	assert.Equal(t, 2, info.DebugId.TimeDateStamp)
	assert.Equal(t, "74219E6FC19941AAB1812BA7E7FF1EFE2", info.DebugId.String())
}

func TestParseExe64(t *testing.T) {
	file, _ := os.Open("../sample/hello64.exe")
	info, err := ParseExe(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "C:\\Bozaro\\pdbuploader\\sample\\hello64.pdb", info.PDBFileName)
	// Code Id
	assert.Equal(t, int32(0x57218B96), info.CodeId.TimeDateStamp)
	assert.Equal(t, int32(0x7e000), info.CodeId.SizeOfImage)
	assert.Equal(t, "57218B967e000", info.CodeId.String())
	// Debug Id
	assert.Equal(t, "B39838386DAE42DF9F1E6B0B12963079", info.DebugId.Guid.String())
	assert.Equal(t, 2, info.DebugId.TimeDateStamp)
	assert.Equal(t, "B39838386DAE42DF9F1E6B0B129630792", info.DebugId.String())
}
