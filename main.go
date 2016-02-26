// https://rsdn.ru/article/baseserv/pe_coff.xml
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

type RvaAndSize struct {
	VirtualAddress int32
	VirtualSize    int32
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

type PESection struct {
	Name                 [0x08]byte
	VirtualSize          int32
	VirtualAddress       int32
	SizeOfRawData        int32
	PointerToRawData     int32
	PointerToRelocations int32
	PointerToLinenumbers int32
	NumberOfRelocations  int16
	NumberOfLinenumbers  int16
	Characteristics      int32
}
type PEDebugDirectory struct {
	Characteristics  int32
	TimeDateStamp    int32
	MajorVersion     int16
	MinorVersion     int16
	Type             int32
	SizeOfData       int32
	AddressOfRawData int32
	PointerToRawData int32
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

	fmt.Printf("Sections count: %d\n", pe.NumberOfSections)
	fmt.Printf("Sections alignment: %d\n", pe.SectionAlignment)
	fmt.Printf("Size of headers: %d\n", pe.SizeOfHeaders)

	var debug_rva RvaAndSize
	if pe.NumberOfRvaAndSizes < 7 {
		// todo: opss...
		fmt.Println("OPS....")
	}
	for i := 0; i < 7; i++ {
		binary.Read(file, binary.LittleEndian, &debug_rva)
	}

	file.Seek(int64(mz.PEOffset)+int64(pe.SizeOfOptionalHeader)+0x18, 0)

	fmt.Printf("Section offset: %08X\n", int64(mz.PEOffset)+int64(pe.SizeOfOptionalHeader)+0x18)
	rdata := [8]byte{'.', 'r', 'd', 'a', 't', 'a'}
	debug_dir_offest := int64(0)
	for i := int16(0); i < pe.NumberOfSections; i++ {
		var section PESection
		binary.Read(file, binary.LittleEndian, &section)
		fmt.Printf("%d: %s\n", i, section.Name)
		if section.Name == rdata {
			debug_dir_offest = int64(section.PointerToRawData + debug_rva.VirtualAddress - section.VirtualAddress)
			break
		}
	}

	if debug_dir_offest > 0 {
		file.Seek(int64(debug_dir_offest), 0)
		var debug_dir PEDebugDirectory
		fmt.Printf("IMAGE_DEBUG_DIRECTORY offset: %X (%d)\n", debug_dir_offest, binary.Size(&debug_dir))
		for i := 0; i < int(debug_rva.VirtualSize)/binary.Size(&debug_dir); i++ {
			binary.Read(file, binary.LittleEndian, &debug_dir)
			fmt.Printf("   %d: %d\n", i, debug_dir.Type)
			if debug_dir.Type == 2 {
				fmt.Printf("RSDS offset: %X\n", debug_dir.PointerToRawData)
				break
			}
		}
	}

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
