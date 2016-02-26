package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type debug_info struct {
	CodeId  string
	DebugId string
}

type MZHeader struct {
	Signature int16      // 0x00-0x02 0x5A4D
	Unused    [0x3A]byte // 0x02-0x3C
	PEOffset  int32      // 0x3C-0x40
}

type PEHeader struct {
	Signature                   int32      // 0x00-0x04 0x50450000
	Machine                     int16      // 0x04-0x06
	NumberOfSections            int16      // 0x06-0x08
	TimeDateStamp               int32      // 0x08-0x0C
	PointerToSymbolTable        int32      // 0x0C-0x10
	NumberOfSymbolTable         int32      // 0x10-0x14
	SizeOfOptionalHeader        int16      // 0x14-0x16
	Characteristics             int16      // 0x16-0x18
	StandadCOFFFields           [0x1C]byte // 0x18-0x34
	ImageBase                   int32      // 0x34-0x38
	SectionAlignment            int32      // 0x38-0x3C
	FileAlignment               int32      // 0x3C-0x40
	MajorOperatingSystemVersion int16      // 0x40-0x42
	MinorOperatingSystemVersion int16      // 0x40-0x42
	MajorImageVersion           int16      // 0x42-0x44
	MinorImageVersion           int16      // 0x44-0x46
	MajorSubsystemVersion       int16      // 0x46-0x48
	MinorSubsystemVersion       int16      // 0x48-0x4A
	Win32VersionValue           int32      // 0x4A-0x50
	SizeOfImage                 int32      // 0x50-0x54
	SizeOfHeaders               int32      // 0x54-0x58
	CheckSum                    int32      // 0x58-0x5C
	Subsystem                   int16      // 0x5C-0x5E
	DllCharacteristics          int16      // 0x5E-0x60
	SizeOfStackReserve          int32      // 0x60-0x64
	SizeOfStackCommit           int32      // 0x64-0x68
	SizeOfHeapReserve           int32      // 0x68-0x6C
	SizeOfHeapCommit            int32      // 0x6C-0x70
	LoaderFlags                 int32      // 0x70-0x74
	NumberOfRvaAndSizes         int32      // 0x74-0x78
}

func read_debug_info(file *os.File) debug_info {
	var mz MZHeader
	var pe PEHeader
	binary.Read(file, binary.LittleEndian, &mz)

	fmt.Printf("MZ signature: %04X\n", mz.Signature)
	fmt.Printf("PE offset: %08X\n", mz.PEOffset)

	file.Seek(int64(mz.PEOffset), 0)
	binary.Read(file, binary.LittleEndian, &pe)

	fmt.Printf("PE signature: %08X\n", pe.Signature)
	fmt.Printf("PE timestamp: %08X\n", pe.TimeDateStamp)
	fmt.Printf("PE image size: %08X\n", pe.SizeOfImage)

	return debug_info{
		fmt.Sprintf("%X%x", pe.TimeDateStamp, pe.SizeOfImage),
		"debug",
	}
}

func main() {
	file, _ := os.Open("sample/hello.exe")
	info := read_debug_info(file)
	fmt.Printf("Code ID: %s\n", info.CodeId)
}
