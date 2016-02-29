package parse

import (
	"fmt"
)

type Guid struct {
	raw [0x10]byte
}

type DebugId struct {
	Guid          Guid
	TimeDateStamp int
}

type CodeId struct {
	TimeDateStamp int32
	SizeOfImage   int32
}

type DebugInfo struct {
	CodeId      CodeId
	DebugId     DebugId
	PDBFileName string
}

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func (guid Guid) String() string {
	raw := guid.raw
	return fmt.Sprintf("%02X%02X%02X%02X%02X%02X%02X%02X%16X",
		raw[3], raw[2], raw[1], raw[0],
		raw[5], raw[4],
		raw[7], raw[6],
		raw[8:])
}

func (this DebugId) String() string {
	return fmt.Sprintf("%s%d", this.Guid, this.TimeDateStamp)
}

func (this CodeId) String() string {
	return fmt.Sprintf("%X%x", this.TimeDateStamp, this.SizeOfImage)
}
