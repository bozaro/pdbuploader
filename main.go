// https://rsdn.ru/article/baseserv/pe_coff.xml
// http://www.godevtool.com/Other/pdb.htm
package pdbuploader

import (
	"fmt"
	"os"
)

type DebugInfo struct {
	CodeId      string
	DebugId     string
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

func guid_to_string(guid [0x10]byte) string {
	return fmt.Sprintf("%02X%02X%02X%02X%02X%02X%02X%02X%16X",
		guid[3], guid[2], guid[1], guid[0],
		guid[5], guid[4],
		guid[7], guid[6],
		guid[8:])
}

func main() {
	file, _ := os.Open("sample/hello.exe")
	{
		info, err := ParseExe(file)
		fmt.Printf("EXE: %s\n", err)
		fmt.Printf("  Code ID: %s\n", info.CodeId)
		fmt.Printf("  Debug ID: %s\n", info.DebugId)
		fmt.Printf("  PDB: %s\n", info.PDBFileName)
	}
	file, _ = os.Open("sample/hello.pdb")
	{
		debug_id, err := ParsePdb(file)
		fmt.Printf("PDB: %s\n", err)
		fmt.Printf("  Debug ID: %s\n", debug_id)
	}
}
